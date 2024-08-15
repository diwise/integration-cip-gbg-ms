package lookup

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
)

type Lookup struct {
	ServiceGuidenId string
	NutsCode        string
	DeviceId        string
}

type LookupTable interface {
	GetNutsCode(serviceGuidenId string) (string, bool)
	GetDeviceId(serviceguidenId string) (string, bool)
}

type impl struct {
	table map[string]*Lookup
}

func fatal(log *slog.Logger, msg string, args ...string) {
	log.Error(fmt.Sprintf(msg, args))
	panic(msg)
}

func New(logger *slog.Logger, filePath string) LookupTable {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fatal(logger, "file %s does not exist", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		fatal(logger, "unable to open file %s", filePath)
	}
	defer file.Close()

	data, err := load(logger, file)
	if err != nil {
		fatal(logger, "unable to load data from file %s", filePath)
	}

	return &impl{
		table: data,
	}
}

func load(_ *slog.Logger, file io.Reader) (map[string]*Lookup, error) {
	r := csv.NewReader(file)
	r.Comma = ';'

	refs, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv data from file: %s", err.Error())
	}

	data := map[string]*Lookup{}

	for idx, r := range refs {
		if idx == 0 {
			continue
		}

		l := &Lookup{
			ServiceGuidenId: r[1],
			NutsCode:        r[2],
			DeviceId:        r[3],
		}

		data[l.ServiceGuidenId] = l
	}

	return data, nil
}

func (l impl) GetNutsCode(serviceGuidenId string) (string, bool) {
	if v, ok := l.table[serviceGuidenId]; ok {
		if v.NutsCode == "" {
			return "", false
		}
		return v.NutsCode, true
	}

	return "", false
}

func (l impl) GetDeviceId(serviceGuidenId string) (string, bool) {
	if v, ok := l.table[serviceGuidenId]; ok {
		return v.DeviceId, true
	}

	return "", false
}
