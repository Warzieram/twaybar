package waybar

import (
	"encoding/json"
	"fmt"
)


type FormatOutput struct{
	Text string `json:"text"`
	Tooltip string `json:"tooltip,omitempty"`
}

func (output FormatOutput) Print() error {
	data, err := json.Marshal(output)	
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
