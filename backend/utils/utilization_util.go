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
