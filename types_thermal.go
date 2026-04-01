package onvif

// ThermalPolarity represents the polarity configuration of a thermal device.
type ThermalPolarity string

const (
	// ThermalPolarityWhiteHot indicates white-hot polarity.
	ThermalPolarityWhiteHot ThermalPolarity = "WhiteHot"
	// ThermalPolarityBlackHot indicates black-hot polarity.
	ThermalPolarityBlackHot ThermalPolarity = "BlackHot"
)

// ThermalColorPaletteType describes standard color palette types for thermal devices.
type ThermalColorPaletteType string

const (
	// ThermalColorPaletteCustom indicates a custom color palette.
	ThermalColorPaletteCustom ThermalColorPaletteType = "Custom"
	// ThermalColorPaletteGrayscale indicates a grayscale color palette.
	ThermalColorPaletteGrayscale ThermalColorPaletteType = "Grayscale"
	// ThermalColorPaletteBlackHot indicates a black-hot color palette.
	ThermalColorPaletteBlackHot ThermalColorPaletteType = "BlackHot"
	// ThermalColorPaletteWhiteHot indicates a white-hot color palette.
	ThermalColorPaletteWhiteHot ThermalColorPaletteType = "WhiteHot"
	// ThermalColorPaletteSepia indicates a sepia color palette.
	ThermalColorPaletteSepia ThermalColorPaletteType = "Sepia"
	// ThermalColorPaletteRed indicates a red color palette.
	ThermalColorPaletteRed ThermalColorPaletteType = "Red"
	// ThermalColorPaletteIron indicates an iron color palette.
	ThermalColorPaletteIron ThermalColorPaletteType = "Iron"
	// ThermalColorPaletteRain indicates a rain color palette.
	ThermalColorPaletteRain ThermalColorPaletteType = "Rain"
	// ThermalColorPaletteRainbow indicates a rainbow color palette.
	ThermalColorPaletteRainbow ThermalColorPaletteType = "Rainbow"
	// ThermalColorPaletteIsotherm indicates an isotherm color palette.
	ThermalColorPaletteIsotherm ThermalColorPaletteType = "Isotherm"
)

// ThermalColorPalette describes a color palette element.
type ThermalColorPalette struct {
	Token string
	Name  string
	Type  string
}

// ThermalNUCTable describes a Non-Uniformity Correction (NUC) table element.
type ThermalNUCTable struct {
	Token           string
	Name            string
	LowTemperature  *float32
	HighTemperature *float32
}

// ThermalCooler describes the cooler settings of a thermal device.
type ThermalCooler struct {
	Enabled bool
	RunTime *float32
}

// ThermalCoolerOptions describes valid ranges for cooler settings.
type ThermalCoolerOptions struct {
	Enabled *bool
}

// ThermalConfiguration holds the thermal settings for a video source.
type ThermalConfiguration struct {
	ColorPalette ThermalColorPalette
	Polarity     ThermalPolarity
	NUCTable     *ThermalNUCTable
	Cooler       *ThermalCooler
}

// ThermalConfigurations holds thermal settings associated with a video source token.
type ThermalConfigurations struct {
	Token         string
	Configuration ThermalConfiguration
}

// ThermalConfigurationOptions describes valid ranges for thermal configuration parameters.
type ThermalConfigurationOptions struct {
	ColorPalettes []*ThermalColorPalette
	NUCTables     []*ThermalNUCTable
	CoolerOptions *ThermalCoolerOptions
}

// ThermalServiceCapabilities represents the capabilities of the thermal service.
type ThermalServiceCapabilities struct {
	Radiometry *bool
}

// RadiometryGlobalParameters holds default values used in radiometry measurement modules.
type RadiometryGlobalParameters struct {
	ReflectedAmbientTemperature float32
	Emissivity                  float32
	DistanceToObject            float32
	RelativeHumidity            *float32
	AtmosphericTemperature      *float32
	AtmosphericTransmittance    *float32
	ExtOpticsTemperature        *float32
	ExtOpticsTransmittance      *float32
}

// RadiometryFloatRange represents a float range for radiometry parameters (uses float32 for ONVIF compatibility).
type RadiometryFloatRange struct {
	Min float32
	Max float32
}

// RadiometryGlobalParameterOptions describes valid ranges for radiometry parameters.
type RadiometryGlobalParameterOptions struct {
	ReflectedAmbientTemperature RadiometryFloatRange
	Emissivity                  RadiometryFloatRange
	DistanceToObject            RadiometryFloatRange
	RelativeHumidity            *RadiometryFloatRange
	AtmosphericTemperature      *RadiometryFloatRange
	AtmosphericTransmittance    *RadiometryFloatRange
	ExtOpticsTemperature        *RadiometryFloatRange
	ExtOpticsTransmittance      *RadiometryFloatRange
}

// RadiometryConfiguration holds the radiometry configuration for a video source.
type RadiometryConfiguration struct {
	RadiometryGlobalParameters *RadiometryGlobalParameters
}

// RadiometryConfigurationOptions describes valid ranges for radiometry configuration parameters.
type RadiometryConfigurationOptions struct {
	RadiometryGlobalParameterOptions *RadiometryGlobalParameterOptions
}
