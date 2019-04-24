package odata

import (
	"encoding/json"
	"errors"
	"io"
)

// JSONReviver is the data type used to parse JSON streams.
type JSONReviver struct {
	decoder *json.Decoder
}

// NewJSONReviver returns a new JSON reviver.
func NewJSONReviver(stream io.Reader) *JSONReviver {
	r := new(JSONReviver)
	r.decoder = json.NewDecoder(stream)

	return r
}

// ParseTransactionLogs parses an incoming stream response that contains transaction log entries.
func (r *JSONReviver) ParseTransactionLogs(callback func(*TransactionLogEntry, bool)) error {
	t, err := r.decoder.Token()
	if err != nil {
		return err
	}

	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return errors.New("JSON object start delimiter not found")
	}

	for r.decoder.More() {
		token, err := r.decoder.Token()
		if err != nil {
			return err
		}

		// Skip other fields than 'value' for simplicity
		if token != "value" {
			continue
		}

		// 'value' should contain an array
		token, err = r.decoder.Token()
		if err != nil {
			return err
		}

		if delim, ok := token.(json.Delim); !ok || delim != '[' {
			return errors.New("JSON array start delimiter not found")
		}

		// Read array elements
		for r.decoder.More() {
			// Read next item (large object)
			txnLog := TransactionLogEntry{}
			err := r.decoder.Decode(&txnLog)
			if err != nil {
				return errors.New("unable to decode transaction log entry")
			}

			// Give transactionLog to the callback for processing.
			callback(&txnLog, false)
		}
		// End of Array
		token, err = r.decoder.Token()
		if err != nil {
			return err
		}

		if delim, ok := token.(json.Delim); !ok || delim != ']' {
			return errors.New("JSON array end delimiter not found")
		}
	}

	// Done parsing
	callback(nil, true)

	return nil
}
