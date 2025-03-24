# mockfactory

**Mockfactory is a simple CLI util to generate mock and test data from your Golang structs.** 

# Installation

***Via go install:***

```bash
go install github.com/maksemen2/mockfactory/cmd/mockfactory@latest
```

***Building from source:***

```bash
git clone https://github.com/maksemen2/mockfactory.git
cd mockfactory
make
```

# Features

- **Realistic random data generation** for your structs fields
- **Flexible data export**:
  - File per struct
  - All data in one file
  - Custom file name templates
- **Field ignore strategies**:
  - Ignore untagged fields
  - Ignore fields with "ignore" tag
  - Generate every field
- **Out-of-box Go primitives support**:
  - Integer (int, uint etc.)
  - Floating (float32, float64)
  - Strings
  - Time

# Examples

**Base**

```go
// test.go
package test

import "time"

type User struct {
	ID        int       `mock:"min=1000;max=9999"`
	Name      string    `mock:"prefix=user"`
	CreatedAt time.Time `mock:"range=past"`
	Balance   float32
}
```

```bash
mockfactory --input ./test.go --count 2 --template "mock_{struct}_{count}" --output result
```

Output in result/mock_User_2.json file:
```json
[
 {
  "Balance": -3.511347e+37,
  "CreatedAt": "2025-03-23T23:34:48.9917166+03:00",
  "ID": 9499,
  "Name": "userkgmRybLT"
 },
 {
  "Balance": 2.4197568e+38,
  "CreatedAt": "2025-03-23T23:34:48.9932214+03:00",
  "ID": 1607,
  "Name": "usermf6RYzcE"
 }
]
```

# Configuration

***CLI arguments***


| Flag | Description | Default |
| ---- | ----------- | ------- |
| --count | Number of objects to generate per struct | 1 |
| --format | Output format | json |
| --ignore | Ignore strategy: untagged or with-tag or all or none | all |
| -i or --input | Path to input Go file | - |
| --log-level | Log level: debug or info or warn or error | error |
| -o or --output | Output path | . |
| --seed | Random seed | time.Now().UnixNano() |
| --strategy | Output strategy: per-struct or single-file | per-struct |
| --structs | Comma-separated list of struct names | - |
| --template | File name template (e.g. {struct}_{count}.json) | - |

**Mock tags**

float(32/64)

| Tag | Description | Default |
| ---- | ----------- | ------- |
| min | Minimal number | -math.MaxFloat32 / -1e307 |
| max | Maximal number | math.MaxFloat32 / 1e307 |

int(-/8/16/32/64)

| Tag | Description | Default |
| ---- | ----------- | ------- |
| min | Minimal number | math.MinInt(-/8/16/32/64) |
| max | Maximal number | math.MaxInt(-/8/16/32/64) |

uint(-/8/16/32/64)

| Tag | Description | Default |
| ---- | ----------- | ------- |
| min | Minimal number | 0 |
| max | Maximal number | math.MaxUint(-/8/16/32/64) |

string

| Tag | Description | Default |
| ---- | ----------- | ------- |
| len | Length of result string (not including prefix and suffix length) | 8 |
| prefix | Prefix of result string | "" |
| suffix | Suffix of result string | "" |

time.Time

| Tag | Description | Default |
| ---- | ----------- | ------- |
| range | Range of the time. Can be "past" or "future" | "" |

github.com/google/uuid.UUID

There are no tags currently available :)

# TBD

- Add nested structures support
- Add more interesting tags for data types (m.b. something like "email" tag for string for email to be created)
- Add convenient API to register your own writers and generators for arbitrary data types