package gowu

type AstronomyResponse struct {
	Response struct {
		Version        string `json:"version"`
		TermsofService string `json:"termsofService"`
		Features       struct {
			Astronomy int `json:"astronomy"`
		} `json:"features"`
	} `json:"response"`
	MoonPhase MoonPhase `json:"moon_phase"`
	SunPhase  SunPhase  `json:"sun_phase"`
}

type HourMinute struct {
	Hour   string `json:"hour"`
	Minute string `json:"minute"`
}

type MoonPhase struct {
	SunRise            HourMinute `json:"sunrise"`
	SunSet             HourMinute `json:"sunset"`
	MoonRise           HourMinute `json:"moonrise"`
	MoonSet            HourMinute `json:"moonset"`
	CurrentTime        HourMinute `json:"current_time"`
	Hemisphere         string     `json:"hemisphere"`
	PercentIlluminated string     `json:"percentIlluminated"`
	AgeOfMoon          string     `json:"ageOfMoon"`
	PhaseOfMoon        string     `json:"phaseofMoon"`
}

type SunPhase struct {
	SunRise HourMinute `json:"sunrise"`
	SunSet  HourMinute `json:"sunset"`
}
