package steamgriddb

type FilterGrid func(d *GridResponse)

func FilterGridStyle(style string) FilterGrid {
	return func(res *GridResponse) {
		var data = []GridResponseData{}
		for _, item := range res.Data {
			if item.Style != style {
				continue
			}
			data = append(data, item)
		}
		res.Data = data
	}
}

type FilterHeroes func(d *HeroesResponse)

func FilterHeroesStyle(style string) FilterHeroes {
	return func(res *HeroesResponse) {
		var data = []ImageResponseData{}
		for _, item := range res.Data {
			if item.Style != style {
				continue
			}
			data = append(data, item)
		}
		res.Data = data
	}
}

type FilterIcons func(d *IconsResponse)

func FilterIconsStyle(style string) FilterIcons {
	return func(res *IconsResponse) {
		var data = []ImageResponseData{}
		for _, item := range res.Data {
			if item.Style != style {
				continue
			}
			data = append(data, item)
		}
		res.Data = data
	}
}

type FilterLogos func(d *LogosResponse)

func FilterLogosStyle(style string) FilterLogos {
	return func(res *LogosResponse) {
		var data = []ImageResponseData{}
		for _, item := range res.Data {
			if item.Style != style {
				continue
			}
			data = append(data, item)
		}
		res.Data = data
	}
}
