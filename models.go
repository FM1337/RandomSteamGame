package main

type SteamResponse struct {
	Response struct {
		Games []struct {
			Appid int `json:"appid"`
		} `json:"games"`
	} `json:"response"`
}

type SteamClientResponse struct {
	Response struct {
		BytesAvailable string `json:"bytes_available"`
		Apps           []struct {
			Appid               int    `json:"appid"`
			App                 string `json:"app"`
			AppType             string `json:"app_type"`
			BytesDownloaded     string `json:"bytes_downloaded"`
			BytesToDownload     string `json:"bytes_to_download"`
			AutoUpdate          bool   `json:"auto_update"`
			Installed           bool   `json:"installed"`
			Changing            bool   `json:"changing"`
			AvailableOnPlatform bool   `json:"available_on_platform"`
			BytesStaged         string `json:"bytes_staged"`
			BytesToStage        string `json:"bytes_to_stage"`
			SourceBuildid       int    `json:"source_buildid"`
			TargetBuildid       int    `json:"target_buildid"`
			QueuePosition       int    `json:"queue_position"`
			Running             bool   `json:"running"`
		} `json:"apps"`
		RefetchIntervalSecFull     int `json:"refetch_interval_sec_full"`
		RefetchIntervalSecChanging int `json:"refetch_interval_sec_changing"`
		RefetchIntervalSecUpdating int `json:"refetch_interval_sec_updating"`
	} `json:"response"`
}

type AuthToken struct {
	Iss         string   `json:"iss"`
	Sub         string   `json:"sub"`
	Aud         []string `json:"aud"`
	Exp         int      `json:"exp"`
	Nbf         int      `json:"nbf"`
	Iat         int      `json:"iat"`
	Jti         string   `json:"jti"`
	Oat         int      `json:"oat"`
	RtExp       int      `json:"rt_exp"`
	Per         int      `json:"per"`
	IpSubject   string   `json:"ip_subject"`
	IpConfirmer string   `json:"ip_confirmer"`
}
