package telemetry

var (
	DefaultObjectives = map[float64]float64{
		0.5:   0.01,
		0.95:  0.001,
		0.99:  0.001,
		0.999: 0.0001,
		1.0:   0,
	}

	DefaultHistogramBuckets = []float64{
		0.001,
		0.01,
		0.1,
		0.2,
		0.3,
		0.4,
		0.45,
		0.5,
		0.55,
		0.6,
		0.65,
		0.7,
		0.75,
		0.8,
		0.85,
		0.9,
		1.0,
		1.5,
		2.0,
		3.0,
		5.0,
		10.0,
		30.0,
		60.0,
		120.0,
		300.0,
	}
)

// ErrLabel is error static label.
const ErrLabel = "error"

// ErrLabelValue returns string representation of error label value.
func ErrLabelValue(err error) string {
	if err != nil {
		return "true"
	}

	return "false"
}
