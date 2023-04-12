package json

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

func Unmarshal(in []byte) ([]ProfileCollection, error) {
	var arrayOfRawJSON []json.RawMessage

	if err := json.Unmarshal(in, &arrayOfRawJSON); err != nil {
		return nil, err
	}

	ret := make([]ProfileCollection, 0)
	var profileCollection *ProfileCollection
	for _, v := range arrayOfRawJSON {
		var genericProfile Profile
		if err := json.Unmarshal(v, &genericProfile); err != nil {
			return nil, err
		}

		switch genericProfile.Profile {
		case "Main":
			if profileCollection != nil {
				ret = append(ret, *profileCollection)
				profileCollection = nil
			}

			var tmp MainProfile
			if err := json.Unmarshal(v, &tmp); err != nil {
				return ret, errors.Wrapf(err, "error unmarshalling next main")
			}

			profileCollection = &ProfileCollection{
				MainProfile: tmp,
			}
		case "security":
			if profileCollection == nil {
				return ret, errors.New("main expected; found security")
			}
			if profileCollection.Sha256 != genericProfile.Sha256 {
				return ret, errors.New(fmt.Sprintf("sha256 mismatch: %s != %s", profileCollection.Sha256, genericProfile.Sha256))
			}

			var securityProfile SecurityProfile
			if err := json.Unmarshal(v, &securityProfile); err != nil {
				return nil, err
			}
			profileCollection.SecurityProfile = &securityProfile

		case "quality":
			if profileCollection == nil {
				return ret, errors.New("main expected; found quality")
			}
			if profileCollection.Sha256 != genericProfile.Sha256 {
				return ret, errors.New(fmt.Sprintf("sha256 mismatch: %s != %s", profileCollection.Sha256, genericProfile.Sha256))
			}

			var qualityProfile QualityProfile
			if err := json.Unmarshal(v, &qualityProfile); err != nil {
				return nil, err
			}
			profileCollection.QualityProfile = &qualityProfile
		default:
			if profileCollection == nil {
				return ret, errors.New(fmt.Sprintf("main expected; found %s", genericProfile.Profile))
			}
			if profileCollection.Sha256 != genericProfile.Sha256 {
				return ret, errors.New(fmt.Sprintf("sha256 mismatch: %s != %s", profileCollection.Sha256, genericProfile.Sha256))
			}

			if profileCollection.UnexpectedProfiles == nil {
				profileCollection.UnexpectedProfiles = map[string]json.RawMessage{genericProfile.Profile: v}
			} else {
				profileCollection.UnexpectedProfiles[genericProfile.Profile] = v
			}
		}
	}
	if profileCollection != nil {
		ret = append(ret, *profileCollection)
	}

	return ret, nil
}
