package fetchrequest

import (
	// stdlib
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	// 3rd-party
	"golang.org/x/net/html"
)

func getChronicle(profilePageTokenizer *html.Tokenizer) string {
	profilePageTokenizer.Next()
	profilePageTokenizer.Next()
	profilePageTokenizer.Next()
	profilePageTokenizer.Next()
	profilePageTokenizer.Next()
	profilePageTokenizer.Next()
	chronicle := ""
	for {
		tt := profilePageTokenizer.Next()
		txt := profilePageTokenizer.Token()
		if tt == html.EndTagToken && txt.Data == "div" {
			break
		} else {
			chronicle += txt.String()
		}
	}

	return chronicle
}

func getSimpleString(profilePageTokenizer *html.Tokenizer, item string, data map[string]string) map[string]string {
	profilePageTokenizer.Next()
	profilePageTokenizer.Next()
	profilePageTokenizer.Next()
	profilePageTokenizer.Next()
	txt := profilePageTokenizer.Token()
	if data[item] == "" {
		data[item] = txt.Data
	}
	return data
}

func getSkill(profilePageTokenizer *html.Tokenizer, data map[string]string) map[string]string {
	numbersRegexp := regexp.MustCompile("[0-9]+")
	profilePageTokenizer.Next()
	txt := profilePageTokenizer.Token()
	possiblySkillName := txt.Data
	profilePageTokenizer.Next()
	txt = profilePageTokenizer.Token()
	profilePageTokenizer.Next()
	txt = profilePageTokenizer.Token()
	if strings.Contains(txt.Data, "уровня") {
		label := strconv.Itoa(int(len(data) / 2))
		labelLvl := strconv.Itoa(int(len(data))/2) + "_level"
		data[label] = possiblySkillName
		data[labelLvl] = strings.Join(numbersRegexp.FindAllString(txt.Data, -1), "")
	}
	return data
}

func getPantheon(profilePageTokenizer *html.Tokenizer, name string, pantheons map[string]string) map[string]string {
	numbersRegexp := regexp.MustCompile("[0-9]+")
	profilePageTokenizer.Next()
	profilePageTokenizer.Next()
	profilePageTokenizer.Next()
	txt := profilePageTokenizer.Token()
	pantheons[name] = strings.Join(numbersRegexp.FindAllString(txt.Data, -1), "")
	return pantheons
}

func getWeapon(profilePageTokenizer *html.Tokenizer, weaponPrefix string, data map[string]string) map[string]string {
	profilePageTokenizer.Next()
	profilePageTokenizer.Next()
	profilePageTokenizer.Next()
	profilePageTokenizer.Next()
	txt := profilePageTokenizer.Token()
	if data[weaponPrefix] == "" {
		data[weaponPrefix] = txt.Data
		profilePageTokenizer.Next()
		profilePageTokenizer.Next()
		profilePageTokenizer.Next()
		profilePageTokenizer.Next()
		txt = profilePageTokenizer.Token()
		data[weaponPrefix+"_stat"] = txt.Data
	}
	return data
}

func parseHTML(data []byte) (map[string]string, map[string]string, map[string]string, bool) {
	c.Log.Info("Profile parsing started...")
	profileHTML := bytes.NewReader(data)
	numbersRegexp := regexp.MustCompile("[0-9]+")

	profilePageTokenizer := html.NewTokenizer(profileHTML)
	profile := make(map[string]string)
	skills := make(map[string]string)
	pantheons := make(map[string]string)
	for {
		tt := profilePageTokenizer.Next()
		switch {
		case tt == html.ErrorToken:
			if profile["god_name"] == "Первая полоса" {
				return profile, skills, pantheons, false
			} else {
				return profile, skills, pantheons, true
			}
		case tt == html.StartTagToken:
			t := profilePageTokenizer.Token()
			switch {
			case t.Data == "h2":
				profilePageTokenizer.Next()
				txt := profilePageTokenizer.Token()
				if profile["god_name"] == "" {
					profile["god_name"] = strings.Join(strings.Fields(txt.Data)[:], " ")
				}
			case t.Data == "h3":
				profilePageTokenizer.Next()
				txt := profilePageTokenizer.Token()
				if profile["hero_name"] == "" {
					profile["hero_name"] = strings.Join(strings.Fields(txt.Data)[:], " ")
				}
			case t.Data == "p":
				for _, attr := range t.Attr {
					if attr.Key == "class" {
						switch {
						case attr.Val == "level":
							profilePageTokenizer.Next()
							txt := profilePageTokenizer.Token()
							if profile["level"] == "" {
								profile["level"] = strings.Split(txt.Data, "-")[0]
							}
							profilePageTokenizer.Next()
							profilePageTokenizer.Next()
							txt = profilePageTokenizer.Token()
							if strings.Contains(txt.Data, "торговец") {
								if profile["shop_level"] == "" {
									profile["shop_level"] = strings.Join(numbersRegexp.FindAllString(txt.Data, -1), "")
								}
							}
						case attr.Val == "motto":
							profilePageTokenizer.Next()
							txt := profilePageTokenizer.Token()
							if profile["motto"] == "" {
								profile["motto"] = strings.Join(strings.Fields(txt.Data)[:], " ")
							}
						}
					}
				}
			case t.Data == "td":
				for _, attr := range t.Attr {
					if attr.Key == "class" {
						switch {
						case attr.Val == "label":
							profilePageTokenizer.Next()
							lbl := profilePageTokenizer.Token()
							switch {
							case lbl.Data == "Возраст":
								profile = getSimpleString(profilePageTokenizer, "age", profile)
							case lbl.Data == "Характер":
								profile = getSimpleString(profilePageTokenizer, "personality", profile)
							case lbl.Data == "Гильдия":
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								txt := profilePageTokenizer.Token()
								if profile["guild"] == "" {
									profile["guild"] = strings.Join(strings.Fields(txt.Data)[:], " ")
								}
							case lbl.Data == "Смертей":
								profile = getSimpleString(profilePageTokenizer, "deaths", profile)
							case lbl.Data == "Побед / Поражений":
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								txt := profilePageTokenizer.Token()
								datas := strings.Join(strings.Fields(txt.Data)[:], " ")
								if profile["arena_wins"] == "" {
									profile["arena_wins"] = strings.Split(string(datas), " / ")[0]
									profile["arena_loss"] = strings.Split(string(datas), " / ")[1]
								}
							case lbl.Data == "Храм достроен":
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								txt := profilePageTokenizer.Token()
								if profile["bricks"] == "" {
									profile["church_done"] = txt.Data
									profile["bricks"] = "1000"
								}
							case lbl.Data == "Кирпичей для храма":
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								txt := profilePageTokenizer.Token()
								if profile["bricks"] == "" {
									profile["church_done"] = ""
									profile["bricks"] = strings.Join(numbersRegexp.FindAllString(txt.Data, -1), "")
								}
							case lbl.Data == "Дерева для ковчега":
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								txt := profilePageTokenizer.Token()
								if profile["woods"] == "" {
									profile["boat_done"] = ""
									profile["woods"] = strings.Join(numbersRegexp.FindAllString(txt.Data, -1), "")
								}
							case lbl.Data == "Ковчег достроен":
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								txt := profilePageTokenizer.Token()
								if profile["woods"] == "" {
									profile["boat_done"] = strings.Split(txt.Data, " ")[0]
									profile["woods"] = strings.Join(numbersRegexp.FindAllString(strings.Split(txt.Data, " ")[1], -1), "")
								}
							case lbl.Data == "Твари по паре":
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								txt := profilePageTokenizer.Token()
								if profile["beasts_male"] == "" {
									profile["beasts_done"] = ""
									profile["beasts_male"] = strings.Join(numbersRegexp.FindAllString(strings.Fields(txt.Data)[0], -1), "")
									profile["beasts_female"] = strings.Join(numbersRegexp.FindAllString(strings.Fields(txt.Data)[1], -1), "")
									profile["beasts_pairs"] = strings.Join(numbersRegexp.FindAllString(strings.Split(strings.Split(txt.Data, " ")[1], "ж")[0], -1), "")
								}
							case lbl.Data == "Твари собраны":
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								txt := profilePageTokenizer.Token()
								if profile["beasts_male"] == "" {
									profile["beasts_done"] = txt.Data
									profile["beasts_male"] = "1000"
									profile["beasts_female"] = "1000"
									profile["beasts_pairs"] = "1000"
								}
							case lbl.Data == "Сбережения":
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								txt := profilePageTokenizer.Token()
								if profile["pension"] == "" {
									profile["pension"] = strings.Join(numbersRegexp.FindAllString(strings.Fields(txt.Data)[0], -1), "")
								}
							case lbl.Data == "Лавка":
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								txt := profilePageTokenizer.Token()
								if profile["pension"] == "" {
									profile["pension"] = "30000"
									profile["shop"] = txt.Data
								}
							case lbl.Data == "Питомец":
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								profilePageTokenizer.Next()
								txt := profilePageTokenizer.Token()
								if profile["pet_name"] == "" {
									profile["pet_type"] = txt.Data
									profilePageTokenizer.Next()
									profilePageTokenizer.Next()
									txt = profilePageTokenizer.Token()
									profile["pet_name"] = strings.Fields(txt.Data)[0]
									if strings.Contains(txt.Data, "уровня") {
										profile["pet_level"] = strings.Join(numbersRegexp.FindAllString(strings.Fields(txt.Data)[1], -1), "")
									}
								}
							case lbl.Data == "Оружие":
								profile = getWeapon(profilePageTokenizer, "weapon", profile)
							case lbl.Data == "Щит":
								profile = getWeapon(profilePageTokenizer, "coat", profile)
							case lbl.Data == "Голова":
								profile = getWeapon(profilePageTokenizer, "head", profile)
							case lbl.Data == "Тело":
								profile = getWeapon(profilePageTokenizer, "body", profile)
							case lbl.Data == "Руки":
								profile = getWeapon(profilePageTokenizer, "hands", profile)
							case lbl.Data == "Ноги":
								profile = getWeapon(profilePageTokenizer, "legs", profile)
							case lbl.Data == "Талисман":
								profile = getWeapon(profilePageTokenizer, "talisman", profile)
							}
						case attr.Val == "name":
							profilePageTokenizer.Next()
							nm := profilePageTokenizer.Token()
							switch {
							case nm.Data == "Благодарности":
								pantheons = getPantheon(profilePageTokenizer, "gratitude", pantheons)
							case nm.Data == "Мощи":
								pantheons = getPantheon(profilePageTokenizer, "might", pantheons)
							case nm.Data == "Храмовничества":
								pantheons = getPantheon(profilePageTokenizer, "templehood", pantheons)
							case nm.Data == "Гладиаторства":
								pantheons = getPantheon(profilePageTokenizer, "gladiatorship", pantheons)
							case nm.Data == "Сказаний":
								pantheons = getPantheon(profilePageTokenizer, "storytelling", pantheons)
							case nm.Data == "Поддержки":
								pantheons = getPantheon(profilePageTokenizer, "donation", pantheons)
							case nm.Data == "Мастерства":
								pantheons = getPantheon(profilePageTokenizer, "mastery", pantheons)
							case nm.Data == "Строительства":
								pantheons = getPantheon(profilePageTokenizer, "construction", pantheons)
							case nm.Data == "Звероводства":
								pantheons = getPantheon(profilePageTokenizer, "taming", pantheons)
							case nm.Data == "Живучести":
								pantheons = getPantheon(profilePageTokenizer, "survival", pantheons)
							case nm.Data == "Зажиточности":
								pantheons = getPantheon(profilePageTokenizer, "savings", pantheons)
							case nm.Data == "Величия":
								pantheons = getPantheon(profilePageTokenizer, "glory", pantheons)
							case nm.Data == "Созидания":
								pantheons = getPantheon(profilePageTokenizer, "creation", pantheons)
							case nm.Data == "Разрушения":
								pantheons = getPantheon(profilePageTokenizer, "destruction", pantheons)
							case nm.Data == "Плотничества":
								pantheons = getPantheon(profilePageTokenizer, "wood", pantheons)
							case nm.Data == "Отлова":
								pantheons = getPantheon(profilePageTokenizer, "pairs", pantheons)
							case nm.Data == "Соперничества":
								pantheons = getPantheon(profilePageTokenizer, "duelers", pantheons)
							case nm.Data == "Солидарности":
								pantheons = getPantheon(profilePageTokenizer, "unity", pantheons)
							case nm.Data == "Влиятельности":
								pantheons = getPantheon(profilePageTokenizer, "popularity", pantheons)
							case nm.Data == "Воинственности":
								pantheons = getPantheon(profilePageTokenizer, "duelery", pantheons)
							}
						}
					}
				}
			case t.Data == "div":
				for _, attr := range t.Attr {
					switch {
					case attr.Key == "class":
						if attr.Val == "guild_status" {
							profilePageTokenizer.Next()
							txt := profilePageTokenizer.Token()
							if profile["guild_status"] == "" {
								profile["guild_status"] = strings.Split(strings.Split(txt.Data, "(")[1], ")")[0]
							}
						}
					case attr.Key == "id":
						if attr.Val == "post_content" {
							profile["chronicle"] = getChronicle(profilePageTokenizer)
						}
					}
				}
			case t.Data == "li":
				skills = getSkill(profilePageTokenizer, skills)
			}
		}
	}
}

func startParsing(godName string) (map[string]string, map[string]string, map[string]string, string) {
	profile := make(map[string]string)
	skills := make(map[string]string)
	pantheons := make(map[string]string)
	site := "https://godville.net"
	client := http.Client{}

	resp, err := client.Get(site + "/gods/" + godName)
	if err != nil {
		c.Log.Errorln(err)
		return profile, skills, pantheons, "error"
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		c.Log.Errorln(err)
		return profile, skills, pantheons, "error"
	}

	profile, skills, pantheons, ok := parseHTML(data)
	if !ok {
		return profile, skills, pantheons, "error"
	}

	return profile, skills, pantheons, "success"
}
