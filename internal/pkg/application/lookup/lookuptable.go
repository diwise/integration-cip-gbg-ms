package lookup

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
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

type lookupTable struct {
	table map[string]*Lookup
}

func New(logger zerolog.Logger, filePath string) LookupTable {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Err(err).Msgf("unable to open file %s", filePath)
		panic(err)
	}
	defer file.Close()

	data, err := loaddata(logger, file)
	if err != nil {
		logger.Err(err).Msgf("unable to load data from file %s", filePath)
		panic(err)
	}

	return &lookupTable{
		table: data,
	}
}

func loaddata(log zerolog.Logger, file io.Reader) (map[string]*Lookup, error) {
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

func (l lookupTable) GetNutsCode(serviceGuidenId string) (string, bool) {
	if v, ok := l.table[serviceGuidenId]; ok {
		if v.NutsCode == "" {
			return "", false
		}
		return v.NutsCode, true
	}

	return "", false
}

func (l lookupTable) GetDeviceId(serviceGuidenId string) (string, bool) {
	if v, ok := l.table[serviceGuidenId]; ok {
		return v.DeviceId, true
	}

	return "", false
}
