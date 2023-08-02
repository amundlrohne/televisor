package utils

import "github.com/amundlrohne/televisor/models"

func SplitUtilization(util models.Utilization, split int) models.Utilization {
	splitFloat := float64(split)
	return models.Utilization{
		Quantile: util.Quantile / splitFloat,
		Mean:     util.Mean / splitFloat,
		Stdev:    util.Stdev / splitFloat,
	}
}

func AddUtilizations(util1 models.Utilization, util2 models.Utilization) models.Utilization {
	return models.Utilization{
		Quantile: util1.Quantile + util2.Quantile,
		Mean:     util1.Mean + util2.Mean,
		Stdev:    util1.Stdev + util2.Stdev,
	}
}
