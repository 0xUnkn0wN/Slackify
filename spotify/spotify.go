package spotify

import (
	"os/exec"
	"runtime"
	"strings"
)

type tokenResponse struct {
	Token string `json:"t"`
}

type SpotifySong struct {
	Item struct {
		Artists string `json:"artists"`
		Name    string `json:"name"`
	} `json:"item"`
}

func GetCurrentTitle() (SpotifySong, error) {
	switch runtime.GOOS {
	case "darwin":
		return getCurrentTitleOSX()
	case "win":
		return getCurrentTitleWin()
	}

	return SpotifySong{}, nil
}

func getCurrentTitleOSX() (SpotifySong, error) {
	song := SpotifySong{}

	output, err := exec.Command("bash", "-c", "osascript -e 'tell application \"Spotify\" to (get artist of current track) & \":\" & (get name of current track)'").Output()
	if err != nil {
		return song, err
	}

	segments := strings.Split(string(output), ":")

	song.Item.Name = strings.TrimSpace(segments[1])
	song.Item.Artists = strings.TrimSpace(segments[0])

	return song, nil
}

func getCurrentTitleWin() (SpotifySong, error) {
	// TODO: rewrite in go -> currently C#
	/*
		process = Process.GetProcessesByName("spotify").Concat(Process.GetProcessesByName("Spotify")).FirstOrDefault(p => p.SessionId == currentSessionId);
		if (process == null) return;

		string title = process.MainWindowTitle;
		if (currentTitle == title) return;

		currentTitle = title;
		title = title.Replace("Spotify", "").Replace("spotify", "").Trim().TrimStart('-').Trim();
	*/
	return SpotifySong{}, nil
}
