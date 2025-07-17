// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package dictionary

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDictionary(t *testing.T) {
	// Ensure the directory exists
	dir := "test_dictionaries"
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatalf("failed to create test dir: %v", err)
	}

	// Create and write the test dictionary file
	path := filepath.Join(dir, "test_dict.json")
	tmp := `{
				"STKPW/TAO*EUP": "ztype",
				"STKPW/TAOEUP": "ztype",
				"STKPW/WAOEU/PWA*BG": "zwieback",
				"STKPW/WAOEU/PWABG": "zwieback",
				"STKPW/WEU/TER": "zwitter",
				"STKPW/WEU/TER/KWROPB": "zwitterion",
				"STKPW/WEURT": "zwitter",
				"STKPWHRAOEF": "disbelief",
				"STKPWHRAEU": "display",
				"STKPWHURB": "zhuzh",
				"STKPWRAOE": "disagree",
				"STKPWRAOEU/PWA*BG": "zwieback",
				"STKPWRAOEU/PWABG": "zwieback",
				"STKPWRAOEG": "disagreeing",
				"STKPWRAEF": "stenography",
				"STKPWRAF": "stenograph",
				"STKPWRAF/*ER": "stenographer",
				"STKPWRAFR": "stenographer",
				"STKPWREU/TER": "zwitter",
				"STKPWREU/TER/KWROPB": "zwitterion",
				"STKPWREURT": "zwitter",
				"STKPWREURT/KWROPB": "zwitterion",
				"STKPWRET": "cigarette",
				"STKPWA/STKPWEU/KEU": "tzatziki",
				"STKPWA/SEU/KEU": "tzatziki",
				"STKPWAO": "zoo",
				"STKPWAO/STKPWAO/HRO/SKWREUBG/KWRAL": "zoological"
    }`
	err = os.WriteFile(path, []byte(tmp), 0644)
	if err != nil {
		t.Fatalf("failed to write test dict: %v", err)
	}
	defer os.Remove(path)

	// Load and test dictionary
	dict, maxOutline, err := LoadDictionaries(dir)
	if err != nil {
		t.Fatalf("failed to load dictionary: %v", err)
	}

	if maxOutline != 5 {
		t.Errorf("expected %q, counted %q", 5, maxOutline)
	}

	got := dict["STKPWHRAEU"]
	want := "display"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
