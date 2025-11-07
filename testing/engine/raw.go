// Copyright 2022 Juan Pablo Tosso and the OWASP Coraza contributors
// SPDX-License-Identifier: Apache-2.0

package engine

import (
	"github.com/corazawaf/coraza/v3/testing/profile"
)

var _ = profile.RegisterProfile(profile.Profile{
	Meta: profile.Meta{
		Author:      "blotus",
		Description: "Test if the body processors work",
		Enabled:     true,
		Name:        "rawrbp.yaml",
	},
	Tests: []profile.Test{
		{
			Title: "rawrbp",
			Stages: []profile.Stage{
				{
					Stage: profile.SubStage{
						Input: profile.StageInput{
							URI:    "/index.php?t1=aaa&t2=bbb&t3=ccc",
							Method: "POST",
							Headers: map[string]string{
								"content-type": "foo",
							},
							Data: `foo=bar&foo2=bar2`,
						},
						Output: profile.ExpectedOutput{
							TriggeredRules:    []int{100, 101, 102, 103},
							NonTriggeredRules: []int{},
						},
					},
				},
			},
		},
	},
	Rules: `
SecRequestBodyAccess On
SecRule REQUEST_HEADERS:content-type "foo" "id: 100, phase:1, pass, log, ctl:requestBodyProcessor=RAW"
SecRule REQBODY_PROCESSOR "RAW" "id: 101,phase:2,log,block"
SecRule REQUEST_BODY "@streq foo=bar&foo2=bar2" "id:102, phase:2,log,block"
SecRule RAW_REQUEST_BODY "@streq foo=bar&foo2=bar2" "id:103, phase:2,log,block"
`,
})
