package modlist

type Modlist struct {
	Archives         []Archive   `json:"Archives"`
	Author           string      `json:"Author"`
	Description      string      `json:"Description"`
	Directives       []Directive `json:"Directives"`
	GameType         string      `json:"GameType"`
	Image            string      `json:"Image"`
	Name             string      `json:"Name"`
	Readme           string      `json:"Readme"`
	WabbajackVersion string      `json:"WabbajackVersion"`
	Website          string      `json:"Website"`
	Community        string      `json:"Community"`
	Version          string      `json:"Version"`
	IsNSFW           bool        `json:"IsNSFW"`
}

type Archive struct {
	Hash  string        `json:"Hash"`
	Meta  string        `json:"Meta"`
	Name  string        `json:"Name"`
	Size  int64         `json:"Size"`
	State *ArchiveState `json:"State,omitempty"`
}

type ArchiveStateType string

const (
	HttpDownloaderType         ArchiveStateType = "HttpDownloader, Wabbajack.Lib"
	NexusDownloaderType        ArchiveStateType = "NexusDownloader, Wabbajack.Lib"
	GameFileSourceType         ArchiveStateType = "GameFileSourceDownloader, Wabbajack.Lib"
	WabbajackCDNDownloaderType ArchiveStateType = "WabbajackCDNDownloader+State, Wabbajack.Lib"
	GoogleDriveDownloaderType  ArchiveStateType = "GoogleDriveDownloader, Wabbajack.Lib"
)

type ArchiveState struct {
	Type        ArchiveStateType `json:"$type"`
	Headers     []Header         `json:"Headers,omitempty"`
	Url         *string          `json:"Url,omitempty"`
	Author      *string          `json:"Author,omitempty"`
	Description *string          `json:"Description,omitempty"`
	FileID      *int             `json:"FileID,omitempty"`
	GameName    *string          `json:"GameName,omitempty"`
	ImageURL    *string          `json:"ImageURL,omitempty"`
	IsNSFW      *bool            `json:"IsNSFW,omitempty"`
	ModID       *int             `json:"ModID,omitempty"`
	Name        *string          `json:"Name,omitempty"`
	Version     *string          `json:"Version,omitempty"`
	Game        *string          `json:"Game,omitempty"`
	GameFile    *string          `json:"GameFile,omitempty"`
	GameVersion *string          `json:"GameVersion,omitempty"`
	Hash        *string          `json:"Hash,omitempty"`
	Id          *string          `json:"Id,omitempty"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DirectiveType string

const (
	RemappedInlineFileType DirectiveType = "RemappedInlineFile"
	FromArchiveType        DirectiveType = "FromArchive"
	CreateBSAType          DirectiveType = "CreateBSA"
	InlineFileType         DirectiveType = "InlineFile"
	PatchedFromArchiveType DirectiveType = "PatchedFromArchive"
)

type Directive struct {
	Type            DirectiveType `json:"$type"`
	Hash            string        `json:"Hash"`
	Size            int64         `json:"Size"`
	SourceDataID    *string       `json:"SourceDataID,omitempty"`
	To              string        `json:"To"`
	ArchiveHashPath []string      `json:"ArchiveHashPath,omitempty"`
	FileStates      []FileState   `json:"FileStates,omitempty"`
	State           *BSAState     `json:"State,omitempty"`
	TempID          *string       `json:"TempID,omitempty"`
	FromHash        *string       `json:"FromHash,omitempty"`
	PatchID         *string       `json:"PatchID,omitempty"`
}

type FileStateType string

const (
	BSAFileStateType FileStateType = "BSAFileState, Compression.BSA"
)

type FileState struct {
	Type            FileStateType `json:"$type"`
	FlipCompression bool          `json:"FlipCompression"`
	Index           int           `json:"Index"`
	Path            string        `json:"Path"`
}

type BSAStateType string

const (
	BSAStateTypeConst BSAStateType = "BSAState, Compression.BSA"
)

type BSAState struct {
	Type         BSAStateType `json:"$type"`
	ArchiveFlags int          `json:"ArchiveFlags"`
	FileFlags    int          `json:"FileFlags"`
	Magic        string       `json:"Magic"`
	Version      int          `json:"Version"`
}
