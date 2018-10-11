package util

import (
	"bytes"
	"strings"
)

/*
GetField will loop through a slack message json and return the fields
*/
func Getfield(attachments []interface{}, buffer *bytes.Buffer) {
	for _, attachment := range attachments {
		attachmentI := attachment.(map[string]interface{})
		fields := attachmentI["fields"].([]interface{})

		//Loop through the fields
		for _, field := range fields {
			fieldI := field.(map[string]interface{})
			buffer.WriteString("*")
			buffer.WriteString(fieldI["title"].(string))
			buffer.WriteString("* ")
			buffer.WriteString(fieldI["value"].(string))
			buffer.WriteString("\n")
		}
	}
}

func EscapeInput(input string) string {
	result := strings.Replace(input, "\"\"", "\\\"\\\"", -1)
	result = strings.Replace(result, "\n", "", -1)
	return result

}
