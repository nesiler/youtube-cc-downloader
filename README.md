Absolutely! Here's a comprehensive `README.md` file for your Go project, structured to be informative and helpful for users and contributors:

**README.md**

# YouTube Subtitle Downloader and CSV Generator

This Go project automates the download of YouTube video subtitles and generates a CSV file containing video details and subtitle content.

## Features

- **Subtitle Download:** Downloads subtitles from specified YouTube channels.
- **Error Logging:** Logs failed downloads to `download_errors.log` for troubleshooting.
- **CSV Generation:** Creates a CSV file (`output.csv`) with video ID, title, channel ID, description, and subtitles.
- **Skipping Existing Downloads:** Avoids re-downloading subtitles that already exist.
- **Customizable:** Control the number of videos to download per channel.
- **Clear Mode:** Removes all downloaded files and directories.

## Prerequisites

- **Go:** Ensure you have Go installed on your system.  You can download it from the official website: [https://golang.org/](https://golang.org/)
- **youtube-dl:** Install `youtube-dl` to fetch video information and subtitles. You can install it using `pip`:

   ```bash
   pip install youtube-dl
   ```

## Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/nesiler/youtube-cc-downloader
   ```
2. Navigate to the project directory:
   ```bash
   cd youtube-cc-downloader
   ```
3. Build the project:
   ```bash
   go build
   ```

## Usage

### 1. Prepare the Channels File (`channels.txt`)

- Create a text file named `channels.txt` in the project's root directory.
- Each line should contain a YouTube channel ID (e.g., `UC...`).

### 2. Run the Downloader

- **Default Mode:** Downloads subtitles from channels in `channels.txt` and generates `output.csv`.
   ```bash
   ./download
   ```

- **Download Mode:** Downloads subtitles only, using `channels.txt` and a specified video count.
   ```bash
   ./download download
   ```
   (This will download 15 videos per channel by default.)

- **Clear Mode:** Removes all downloaded files and directories.
   ```bash
   ./download clear
   ```

- **CSV Mode:** Generates `output.csv` from existing subtitle files.
   ```bash
   ./download csv
   ```

## Configuration

- **Number of Videos:** By default, the downloader fetches the latest 15 videos per channel. You can change this by modifying the `count` variable in `main.go`.


## Contributing

Contributions are welcome! If you'd like to enhance this project, please fork the repository and submit a pull request.

## License

This project is licensed under the [GNU License](LICENSE).
