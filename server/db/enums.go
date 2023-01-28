package db

type FuelUnit int

const (
	LITRE FuelUnit = iota
	GALLON
	US_GALLON
	KILOGRAM
	KILOWATT_HOUR
	MINUTE
)

type FuelType int

const (
	PETROL FuelType = iota
	DIESEL
	ETHANOL
	CNG
	ELECTRIC
	LPG
)

type DistanceUnit int

const (
	MILES DistanceUnit = iota
	KILOMETERS
)

type Role int

const (
	ADMIN Role = iota
	USER
)

type AlertFrequency int

const (
	ONETIME AlertFrequency = iota
	RECURRING
)

type AlertType int

const (
	DISTANCE AlertType = iota
	TIME
	BOTH
)

type EnumDetail struct {
	Key string `json:"key"`
}

var FuelUnitDetails map[FuelUnit]EnumDetail = map[FuelUnit]EnumDetail{
	LITRE: {
		Key:  "litre",
	},
	GALLON: {
		Key:  "gallon",
	}, KILOGRAM: {
		Key:  "kilogram",
	}, KILOWATT_HOUR: {
		Key:  "kilowatthour",
	}, US_GALLON: {
		Key:  "usgallon",
	},
	MINUTE: {
		Key:  "minutes",
	},
}

var FuelTypeDetails map[FuelType]EnumDetail = map[FuelType]EnumDetail{
	PETROL: {
		Key:  "petrol",
	},
	DIESEL: {
		Key:  "diesel",
	}, CNG: {
		Key:  "cng",
	}, LPG: {
		Key:  "lpg",
	}, ELECTRIC: {
		Key:  "electric",
	}, ETHANOL: {
		Key:  "ethanol",
	},
}

var DistanceUnitDetails map[DistanceUnit]EnumDetail = map[DistanceUnit]EnumDetail{
	KILOMETERS: {
		Key: "kilometers",
	},
	MILES: {
		Key: "miles",
	},
}

var RoleDetails map[Role]EnumDetail = map[Role]EnumDetail{
	ADMIN: {
		Key: "ADMIN",
	},
	USER: {
		Key: "USER",
	},
}
