package query

//type ast struct {
//}
//
//func ParseTime() {
//
//}

// test for ParseQuery
func ParseQuery(query string) map[string]interface{} {
	return map[string]interface{}{
		"source": "kubernetes",
		"_or": map[string]interface{}{
			"pod_name =":  "pod-2",
			"pod_name !=": "pod-3",
		},
	}
}
