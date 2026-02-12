package worker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ParsedLine struct {
	UnitGuid string
	MsgID    string
	Payload  map[string]any
	LineNo   int
}

type ParseError struct {
	LineNo   int
	UnitGuid *string
	Err      error
	RawLine  string
}

func ParseTSVFile(filePath string) ([]ParsedLine, []ParseError, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	var (
		out     []ParsedLine
		errs    []ParseError
		headers []string
	)

	sc := bufio.NewScanner(f)
	buf := make([]byte, 0, 1024*1024)
	sc.Buffer(buf, 10*1024*1024)

	lineNo := 0
	for sc.Scan() {
		lineNo++
		raw := sc.Text()
		if strings.TrimSpace(raw) == "" {
			continue
		}

		cols := strings.Split(raw, "\t")
		if len(cols) < 4 {
			errs = append(errs, ParseError{LineNo: lineNo, Err: fmt.Errorf("too few columns"), RawLine: raw})
			continue
		}

		unitCell := strings.TrimSpace(cols[3])
		lc := strings.ToLower(unitCell)

		if unitCell == "" || lc == "гуид" || lc == "unit_guid" {
			if len(cols) > 1 {
				// headers = B.. (пропускаем A/#номер)
				headers = make([]string, 0, len(cols)-1)
				for _, h := range cols[1:] {
					headers = append(headers, strings.TrimSpace(h))
				}
			}
			continue
		}

		unitGuid := unitCell

		// E колонка = msg_id
		var msgID string
		if len(cols) > 4 {
			msgID = strings.TrimSpace(cols[4])
		}

		values := cols[1:] // B..
		payload := map[string]any{}

		if len(headers) == 0 {
			payload["cols"] = values
		} else {
			n := len(headers)
			if len(values) < n {
				n = len(values)
			}
			for i := 0; i < n; i++ {
				key := headers[i]
				if key == "" {
					key = fmt.Sprintf("col_%d", i+1)
				}
				payload[key] = strings.TrimSpace(values[i])
			}
			for i := n; i < len(values); i++ {
				payload[fmt.Sprintf("extra_%d", i-n+1)] = strings.TrimSpace(values[i])
			}
		}

		out = append(out, ParsedLine{
			UnitGuid: unitGuid,
			MsgID:    msgID,
			Payload:  payload,
			LineNo:   lineNo,
		})
	}

	if err := sc.Err(); err != nil {
		return nil, errs, err
	}

	for i := range out {
		if _, err := json.Marshal(out[i].Payload); err != nil {
			u := out[i].UnitGuid
			errs = append(errs, ParseError{
				LineNo:   out[i].LineNo,
				UnitGuid: &u,
				Err:      fmt.Errorf("payload json marshal: %w", err),
			})
		}
	}

	return out, errs, nil
}
