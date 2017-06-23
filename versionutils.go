package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"io/ioutil"
	"fmt"
	"os"
)

const (
	apiLinkLibrary		string = "https://registry.hub.docker.com/v2/repositories/library/%s/tags/"
	apiLinkLibraryPage	string = "https://registry.hub.docker.com/v2/repositories/library/%s/tags/?Page=%d"
	apiLinkRepo		string = "https://registry.hub.docker.com/v2/repositories/%s/tags/"
	apiLinkRepoPage		string = "https://registry.hub.docker.com/v2/repositories/%s/tags/?page=%d"

	dt_PageNotFound		string = "Not found"
	dt_ImageNotFound	string = "Object not found"

	info_TypeDefault	string = "default"
	info_TypeCurrent	string = "current"
	info_TypeRollback	string = "rollback"
	info_TypeUpdate		string = "update"
)

//--------------------------------------------------------------
type Info struct {
	Type		string		`json:"type"`
	Name		string		`json:"name"`
	Tag		string		`json:"tag"`
	Id		int64		`json:"id"`
//	Last_updated	string		`json:"last_updated"`
}

//--------------------------------------------------------------
type TagResult struct {
	Name		string 		`json:"name"`
	Id		int64		`json:"id"`
	Last_updated	string		`json:"last_updated"`
}

type Tags struct {
	Detail		string		`json:"detail"`
	Count		int64		`json:"count"`
	Next		string		`json:"next"`
	Previous	string		`json:"previous"`
	Results		[]TagResult	`json:"results"`
}

//--------------------------------------------------------------

//	b: false - repo (images)
//	b: true  - full link 
func getTags(b bool, repo string) (ts Tags, err error) {
	link := fmt.Sprintf(apiLinkRepo, repo)
	if b {
		link = repo
	}
	
	// Get data using api registry.hub.docker/v2
	resp, err := http.Get(link)
	if err != nil {
		return ts, err
	}
	
	// Close connection
	defer resp.Body.Close()
	
	// Read data
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ts, err
	}

	// Unmarshal json
	json.Unmarshal(raw, &ts)
	return ts, nil
}

func trueRepo(image string) string {
	if strings.ContainsAny(image, "/") {
		return image
	}
	return "library/" + image
}

func getAllTags(repo string) (ats []Tags, err error) {
	ts, err := getTags(false, repo)
	if err != nil {
		return ats, err
	}

	ats = append(ats, ts)
	for ts.Next != "" {
		ts, err = getTags(true, ts.Next)
		// Never use ts, err := getTags(true, ts.Next)
		// it will be loop forever because ts.Next unchange
		if err != nil {
			return ats, err
		}
		ats = append(ats, ts)
	}
	return ats, err
}

func findTagName(ats []Tags, name string) bool {
	for i := 0; i < len(ats); i++ {
		for j := 0; j < len(ats[i].Results); j++ {
			if ats[i].Results[j].Name == name {
				return true
			}
		}
	}
	return false
}

func findTagId(ats []Tags, id int64) bool {
	for i := 0; i < len(ats); i++ {
		for j := 0; j < len(ats[i].Results); j++ {
			if ats[i].Results[j].Id == id {
				return true
			}
		}
	}
	return false
}

func findTag(ats []Tags, name string, id int64) (b1 bool, b2 bool) {
	for i := 0; i < len(ats); i++ {
		for j := 0; j < len(ats[i].Results); j++ {
			if ats[i].Results[j].Name == name {
				if ats[i].Results[j].Id == id {
					return true, true
				}
				return true, false
			}
		}
	}
	return false, false
}

func getTagId(ats []Tags, tag string) int64 {
        for i := 0; i < len(ats); i++ {
                for j := 0; j < len(ats[i].Results); j++ {
                        if ats[i].Results[j].Name == tag {
                                return ats[i].Results[j].Id
                        }
                }
        }
        return 0
}

//--------------------------------------------------------------
func getInfo(fn string) (ain []Info, err error) {
	raw, err := ioutil.ReadFile(fn)
	if err != nil {
		return ain, err
	}

	err = json.Unmarshal(raw, &ain)
	return ain, err
}

func getDefault(ain []Info) (in Info, b bool) {
	for i := 0; i < len(ain); i++ {
		if ain[i].Type == info_TypeDefault {
			return ain[i], true
		}
	}
	return in, false
}

func getCurrent(ain []Info) (in Info, b bool) {
        for i := 0; i < len(ain); i++ {
                if ain[i].Type == info_TypeCurrent {
                        return ain[i], true
                }
        }
        return in, false
}

func getRollback(ain []Info) (in Info, b bool) {
	for i := 0; i < len(ain); i++ {
                if ain[i].Type == info_TypeRollback {
                        return ain[i], true
                }
        }       
        return in, false
}

func getUpdate(ain []Info) (in Info, b bool) {
	for i := 0; i < len(ain); i++ {
		if ain[i].Type == info_TypeUpdate {
			return ain[i], true
		}
	}
	return in, false
}

func applyDefault(ain []Info, df Info) []Info {
	df.Type = info_TypeDefault
	for i := 0; i < len(ain); i++ {
		if ain[i].Type == info_TypeDefault {
			ain[i] = df
			return ain
		}
	}
	ain = append(ain, df)
	return ain
}

func applyCurrent(ain []Info, cr Info) []Info {
	cr.Type = info_TypeCurrent
	for i := 0; i < len(ain); i++ {
		if ain[i].Type == info_TypeCurrent {
			ain[i] = cr
			return ain
		}
	}
	ain = append(ain, cr)
	return ain
}

func applyRollback(ain []Info, rb Info) []Info {
	rb.Type = info_TypeRollback
        for i := 0; i < len(ain); i++ {
                if ain[i].Type == info_TypeRollback {
                        ain[i] = rb
			return ain
                }       
        }
	ain = append(ain, rb)
	return ain
}

func applyUpdate(ain []Info, ud Info) []Info {
	for i := 0; i < len(ain); i++ { 
                if ain[i].Type == info_TypeUpdate {
                        ain[i] = ud
                        ain[i].Type = info_TypeUpdate
                        return ain
                }
        }
        ain = append(ain, ud)
	return ain
}

func saveInfo(fn string, ain []Info) error {
	// Convert to Json
	bytes, err := json.MarshalIndent(&ain, "", "\t")
	if err != nil {
		return err
	}

	// Write file
	if err = ioutil.WriteFile(fn, bytes, 0644); err != nil {
		return err
	}
	return nil
}

func subCCF_getName(path string, sep string) string {
	s_split := strings.Split(path, sep)
	// Always s_split have at least an member.
	return s_split[len(s_split) - 1]
}

// b: true - all file good | other some file get error
// return output: list file error
// mark: separate between file
func checkConfigFile(lFile []string, mark string) (output string, b bool) {
	var rs []string
	b = true
	for i := 0; i < len(lFile); i++ {
		fInfo, err := os.Stat(lFile[i])
		if err != nil {
			rs = append(rs, subCCF_getName(lFile[i], "/"))
		} else {
			if fInfo.IsDir() {
				rs = append(rs, fInfo.Name())
			}
		}
	}
	if len(rs) > 0 {
		output += rs[0]
		b = false
	}
	for i := 1; i < len(rs); i++ {
		output += (mark + rs[i])
	}
	return output, b
}
