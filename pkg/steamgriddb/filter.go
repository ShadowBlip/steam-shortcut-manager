package steamgriddb

// FilterGrid is a function signature for any function that will filter grid
// results.
type FilterGrid func(d *GridResponse) []GridResponseData

// FilterGridStyle will return a filter that will filter out all results
// that don't match the given style.
func FilterGridStyle(style string) FilterGrid {
	return func(res *GridResponse) []GridResponseData {
		var data = []GridResponseData{}
		for _, item := range res.Data {
			if item.Style != style {
				continue
			}
			data = append(data, item)
		}
		return data
	}
}

// FilterGridVertical will return a filter that will filter out all results
// that are not vertical poster images.
func FilterGridVertical() FilterGrid {
	return func(res *GridResponse) []GridResponseData {
		var data = []GridResponseData{}
		for _, item := range res.Data {
			wantRatio := float64(600) / float64(900)
			itemRatio := float64(item.Width) / float64(item.Height)
			if wantRatio != itemRatio {
				continue
			}
			data = append(data, item)
		}
		return data
	}
}

// FilterGridHorizontal will return a filter that will filter out all results
// that are not horizontal banner images.
func FilterGridHorizontal() FilterGrid {
	return func(res *GridResponse) []GridResponseData {
		var data = []GridResponseData{}
		for _, item := range res.Data {
			wantRatio := float64(920) / float64(430)
			itemRatio := float64(item.Width) / float64(item.Height)
			if wantRatio != itemRatio {
				continue
			}
			data = append(data, item)
		}
		return data
	}
}

type FilterHeroes func(d *HeroesResponse) []ImageResponseData

func FilterHeroesStyle(style string) FilterHeroes {
	return func(res *HeroesResponse) []ImageResponseData {
		var data = []ImageResponseData{}
		for _, item := range res.Data {
			if item.Style != style {
				continue
			}
			data = append(data, item)
		}
		return data
	}
}

type FilterIcons func(d *IconsResponse) []ImageResponseData

func FilterIconsStyle(style string) FilterIcons {
	return func(res *IconsResponse) []ImageResponseData {
		var data = []ImageResponseData{}
		for _, item := range res.Data {
			if item.Style != style {
				continue
			}
			data = append(data, item)
		}
		return data
	}
}

type FilterLogos func(d *LogosResponse) []ImageResponseData

func FilterLogosStyle(style string) FilterLogos {
	return func(res *LogosResponse) []ImageResponseData {
		var data = []ImageResponseData{}
		for _, item := range res.Data {
			if item.Style != style {
				continue
			}
			data = append(data, item)
		}
		return data
	}
}
