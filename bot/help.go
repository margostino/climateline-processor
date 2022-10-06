package bot

func Help() string {
	return GetUpdateHelp()
}

func GetUpdateHelp() string {
	return "❓ What do you want to do?\n" +
		"➡️ edit {id}\n" +
		"{new title}\n" +
		"{source name}\n" +
		"{location}\n" +
		"{category[agreements | assessment | awareness | warming | wildfires | floods | drought | health | hurricane]}\n" +
		"⚡️️ Example:\n" +
		"edit 1\n" +
		"Massive heatwaves in Europe\n" +
		"Washington Post\n" +
		"Europe\n" +
		"warming\n" +
		"➡️ push {ids}\n" +
		"⚡️️ Example:\n" +
		"push 1 2"
}
