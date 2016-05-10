package clair

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getSampleAnalysis() []byte {
	file, err := ioutil.ReadFile("./samples/clair_report.json")

	if err != nil {
		fmt.Printf("File error: %v\n", err)
	}

	return file
}

func newServer(httpStatus int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(httpStatus)
	}))
}

func TestIsHealthy(t *testing.T) {
	server := newServer(http.StatusOK)
	defer server.Close()
	uri = server.URL
	if h := IsHealthy(); !h {
		t.Errorf("IsHealthy() => %v, want %v", h, true)
	}
}

func TestIsNotHealthy(t *testing.T) {
	server := newServer(http.StatusInternalServerError)
	defer server.Close()
	uri = server.URL
	if h := IsHealthy(); h {
		t.Errorf("IsHealthy() => %v, want %v", h, true)
	}
}

// func TestCountAllVulnerabilities(t *testing.T)  {

//   var analysis ImageAnalysis
//   err := json.Unmarshal([]byte(getSampleAnalysis()), &analysis)

//   if err != nil {
// 		t.Errorf("Failing with error: %v", err)
// 	}

//   vulnerabilitiesCount := analysis.CountAllVulnerabilities()

//   if vulnerabilitiesCount.TotalFeatures != 126 {
//     t.Errorf("analysis.CountAllVulnerabilities().TotalFeatures => %v, want 126", vulnerabilitiesCount.TotalFeatures)
//   }

//   if vulnerabilitiesCount.SafeFeatures != 101 {
//     t.Errorf("analysis.CountAllVulnerabilities().SafeFeatures => %v, want 49", vulnerabilitiesCount.SafeFeatures)
//   }

//   if vulnerabilitiesCount.UnsafeFeatures != 25 {
//     t.Errorf("analysis.CountAllVulnerabilities().UnsafeFeatures => %v, want 25", vulnerabilitiesCount.UnsafeFeatures)
//   }

//   if vulnerabilitiesCount.Total != 77 {
//     t.Errorf("analysis.CountAllVulnerabilities().Total => %v, want 77", vulnerabilitiesCount.Total)
//   }

//   if vulnerabilitiesCount.High != 1 {
//     t.Errorf("analysis.CountAllVulnerabilities().High => %v, want 1", vulnerabilitiesCount.High)
//   }

//   if vulnerabilitiesCount.Medium != 18 {
//     t.Errorf("analysis.CountAllVulnerabilities().Medium => %v, want 18", vulnerabilitiesCount.Medium)
//   }

//   if vulnerabilitiesCount.Low != 57 {
//     t.Errorf("analysis.CountAllVulnerabilities().Low => %v, want 57", vulnerabilitiesCount.Low)
//   }
//   if vulnerabilitiesCount.Negligible != 1 {
//     t.Errorf("analysis.CountAllVulnerabilities().Negligible => %v, want 1", vulnerabilitiesCount.Negligible)
//   }
// }

// func TestRelativeCount(t *testing.T) {
// 	var analysis ImageAnalysis
// 	err := json.Unmarshal([]byte(getSampleAnalysis()), &analysis)

// 	if err != nil {
// 		t.Errorf("Failing with error: %v", err)
// 	}

// 	vulnerabilitiesCount := analysis.CountAllVulnerabilities()

// 	if vulnerabilitiesCount.RelativeCount("High", false) != 1.3 {
// 		t.Errorf("analysis.CountAllVulnerabilities().RelativeCount(\"High\") => %v, want 1.3", vulnerabilitiesCount.RelativeCount("High", false))
// 	}

// 	if vulnerabilitiesCount.RelativeCount("High", true) != 0.26 {
// 		t.Errorf("analysis.CountAllVulnerabilities().RelativeCount(\"High\", true) => %v, want 0.26", vulnerabilitiesCount.RelativeCount("High", true))
// 	}

// 	if vulnerabilitiesCount.RelativeCount("Medium", false) != 23.38 {
// 		t.Errorf("analysis.CountAllVulnerabilities().RelativeCount(\"Medium\") => %v, want 23.38", vulnerabilitiesCount.RelativeCount("Medium", false))
// 	}

// 	if vulnerabilitiesCount.RelativeCount("Medium", true) != 4.64 {
// 		t.Errorf("analysis.CountAllVulnerabilities().RelativeCount(\"Medium\", true) => %v, want 4.64", vulnerabilitiesCount.RelativeCount("Medium", true))
// 	}

// 	if vulnerabilitiesCount.RelativeCount("Low", false) != 74.03 {
// 		t.Errorf("analysis.CountAllVulnerabilities().RelativeCount(\"Low\") => %v, want 74.03", vulnerabilitiesCount.RelativeCount("Low", false))
// 	}

// 	if vulnerabilitiesCount.RelativeCount("Low", true) != 14.69 {
// 		t.Errorf("analysis.CountAllVulnerabilities().RelativeCount(\"Low\", true) => %v, want 14.69", vulnerabilitiesCount.RelativeCount("Low", true))
// 	}
// }

// func TestFeatureWeight(t *testing.T) {
// 	feature := Feature{
// 		Vulnerabilities: []Vulnerability{},
// 	}

// 	v1 := Vulnerability{
// 		Severity: "High",
// 	}

// 	v2 := Vulnerability{
// 		Severity: "Medium",
// 	}

// 	v3 := Vulnerability{
// 		Severity: "Low",
// 	}

// 	v4 := Vulnerability{
// 		Severity: "Negligible",
// 	}

// 	feature.Vulnerabilities = append(feature.Vulnerabilities, v1)
// 	feature.Vulnerabilities = append(feature.Vulnerabilities, v2)
// 	feature.Vulnerabilities = append(feature.Vulnerabilities, v3)
// 	feature.Vulnerabilities = append(feature.Vulnerabilities, v4)

// 	if feature.Weight() != 10 {
// 		t.Errorf("feature.Weigh => %v, want 6", feature.Weight())
// 	}
// }
