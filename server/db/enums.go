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

type EnumDetail struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

var FuelUnitDetails map[FuelUnit]EnumDetail = map[FuelUnit]EnumDetail{
	LITRE: {
		Short: "Lt",
		Long:  "Litre",
	},
	GALLON: {
		Short: "Gal",
		Long:  "Gallon",
	}, KILOGRAM: {
		Short: "Kg",
		Long:  "Kilogram",
	}, KILOWATT_HOUR: {
		Short: "KwH",
		Long:  "Kilowatt Hour",
	}, US_GALLON: {
		Short: "US Gal",
		Long:  "US Gallon",
	},
	MINUTE: {
		Short: "Mins",
		Long:  "Minutes",
	},
}

var FuelTypeDetails map[FuelType]EnumDetail = map[FuelType]EnumDetail{
	PETROL: {
		Short: "Petrol",
		Long:  "Petrol",
	},
	DIESEL: {
		Short: "Diesel",
		Long:  "Diesel",
	}, CNG: {
		Short: "CNG",
		Long:  "CNG",
	}, LPG: {
		Short: "LPG",
		Long:  "LPG",
	}, ELECTRIC: {
		Short: "Electric",
		Long:  "Electric",
	}, ETHANOL: {
		Short: "Ethanol",
		Long:  "Ethanol",
	},
}

var DistanceUnitDetails map[DistanceUnit]EnumDetail = map[DistanceUnit]EnumDetail{
	KILOMETERS: {
		Short: "Km",
		Long:  "Kilometers",
	},
	MILES: {
		Short: "Mi",
		Long:  "Miles",
	},
}

var RoleDetails map[Role]EnumDetail = map[Role]EnumDetail{
	ADMIN: {
		Short: "Admin",
		Long:  "ADMIN",
	},
	USER: {
		Short: "User",
		Long:  "USER",
	},
}
