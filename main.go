package main

import (
    "crypto/sha256"
     "github.com/ujuojadi/PersonalFileS/p2p"
	 "log"
	 "strings"
	 "fmt"
    "net/http"
	"time"
    "io"
    "os"
     "encoding/json" 
   
    "path/filepath"
      
      "mime"
      "strconv"
    "github.com/google/uuid"
   
)

	 
func enableCORS(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        h.ServeHTTP(w, r)
    })
}

	

func makeServer(listenAddr string, nodes ...string ) *FileServer {
	tcptransportOpts :=p2p.TCPTransportOpts {
		ListenAddr: listenAddr,
		HandshakeFunc :p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
		
	}
	tcpTransport :=p2p.NewTCPTransport(tcptransportOpts)
    

	safeRoot := strings.ReplaceAll(listenAddr, ":", "_")
	fileServerOpts :=ServerOpts{
		StorageRoot : safeRoot + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport: tcpTransport, 
		BootstrapNodes: nodes,

	}
	s:= NewFileServer(fileServerOpts)
    tcpTransport.OnPeer = s.OnPeer

return s
 
	
}



// Add this function
// func uploadHandler(s *FileServer) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
//         file, header, err := r.FormFile("file")
//         if err != nil {
//             http.Error(w, err.Error(), http.StatusBadRequest)
//             return
//         }
//         defer file.Close()

//         // Store in P2P network
//         err = s.Store(header.Filename, file)
//         if err != nil {
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//             return
//         }

//         w.Write([]byte("uploaded"))
//     }
// }


// func uploadHandler(s *FileServer) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
//         file, header, err := r.FormFile("file")
//         if err != nil {
//             http.Error(w, err.Error(), http.StatusBadRequest)
//             return
//         }
//         defer file.Close()

//         // Store in P2P network
//         err = s.Store(header.Filename, file)
//         if err != nil {
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//             return
//         }

//         // Create manifest
//         meta := FileMetadata{
//             Hash:       s.ComputeHash(header.Filename), // or whatever your CAS hash function is
//             Filename:   header.Filename,
//             Size:       header.Size,
//             Uploader:   "user123",             // can get from auth/session
//             Course:     "CS101",               // optional
//             UploadedAt: time.Now().Format(time.RFC3339),
//         }

//         // Ensure manifest folder exists
//         os.MkdirAll("manifests", os.ModePerm)
//         manifestPath := fmt.Sprintf("manifests/%s.json", meta.Filename)
//         data, _ := json.Marshal(meta)
//         os.WriteFile(manifestPath, data, 0644)

//         w.Write([]byte("uploaded"))
//     }
// }
func uploadHandler(s *FileServer) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        file, header, err := r.FormFile("file")
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        defer file.Close()

        // Compute SHA256 hash for CAS
        file.Seek(0, io.SeekStart) // rewind
        hasher := sha256.New()
        io.Copy(hasher, file)
        hash := fmt.Sprintf("%x", hasher.Sum(nil))

        file.Seek(0, io.SeekStart) // rewind again before storing
        err = s.Store(header.Filename, file)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Create manifest
        meta := FileMetadata{
            Hash:       hash,
            Filename:   header.Filename,
            Size:       header.Size,
            Uploader:   "user123",
            Course:     "CS101",
            UploadedAt: time.Now().Format(time.RFC3339),
        }

        os.MkdirAll("manifests", os.ModePerm)
        manifestPath := fmt.Sprintf("manifests/%s.json", meta.Filename)
        data, _ := json.Marshal(meta)
        os.WriteFile(manifestPath, data, 0644)

        w.Write([]byte("uploaded"))
    }
}
// File metadata struct
type FileMetadata struct {
    Hash       string `json:"hash"`
    Filename   string `json:"filename"`
    Size       int64  `json:"size"`
    Uploader   string `json:"uploader"`
    Course     string `json:"course"`
    UploadedAt string `json:"uploaded_at"`
}

// List files handler
type FileResponse struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Size        int64  `json:"size"`
    CourseCode  string `json:"courseCode"`
    CourseName  string `json:"courseName"`
    Description string `json:"description"`
    Type        string `json:"type"`
    UploadedAt  string `json:"uploadedAt"`
}

// func listFilesHandler(w http.ResponseWriter, r *http.Request) {
//     manifests := []FileMetadata{}
//     files, _ := os.ReadDir("manifests")
//     for _, file := range files {
//         data, _ := os.ReadFile("manifests/" + file.Name())
//         var meta FileMetadata
//         json.Unmarshal(data, &meta)
//         manifests = append(manifests, meta)
//     }

//     resp := []FileResponse{}
//     for _, m := range manifests {
//         resp = append(resp, FileResponse{
//             ID:          m.Hash,
//             Name:        m.Filename,
//             Size:        m.Size,
//             CourseCode:  m.Course,
//             CourseName:  "",            // optional
//             Description: "",            // optional
//             Type:        "application/octet-stream",
//             UploadedAt:  m.UploadedAt,
//         })
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(resp)
// }
// type File struct {
//     ID          string `json:"ID"`
//     Name        string `json:"Name"`
//     Size        int64  `json:"Size"`
//     Type        string `json:"Type"`
//     CourseCode  string `json:"CourseCode,omitempty"`
//     CourseName  string `json:"CourseName,omitempty"`
//     Description string `json:"description,omitempty"`
//     UploadedAt  int64  `json:"uploadedAt"`
// }

// // filesHandler lists all files in memory or in storage
// func filesHandler(w http.ResponseWriter, r *http.Request) {
//     w.Header().Set("Content-Type", "application/json")
//     files := []File{
//         {ID: "1", Name: "example.pdf", Size: 12345, Type: "application/pdf"},
//         // add more files or load from storage
//     }
//     json.NewEncoder(w).Encode(files)
// }






// type File struct {
//     ID          string `json:"ID"`
//     Name        string `json:"Name"`
//     Size        int64  `json:"Size"`
//     Type        string `json:"Type"`
//     CourseCode  string `json:"CourseCode,omitempty"`
//     CourseName  string `json:"CourseName,omitempty"`
//     Description string `json:"description,omitempty"`
//     UploadedAt  int64  `json:"uploadedAt"`
// }


// func detectMimeType(filename string) string {
//     ext := filepath.Ext(filename)
//     mimeType := mime.TypeByExtension(ext)
//     if mimeType != "" {
//         return mimeType
//     }
//     return "application/octet-stream"
// }

// // Recursively scan a network folder

// func scanNetworkFolder(basePath string) ([]File, error) {
//     var files []File
//     err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
//         if err != nil {
//             return nil // skip inaccessible files
//         }
//         if info.IsDir() {
//             return nil
//         }

//         f := File{
//             ID:         info.Name(),
//             Name:       info.Name(),
//             Size:       info.Size(),
//             Type:       detectMimeType(info.Name()),
//             UploadedAt: info.ModTime().Unix(),
//         }
//         files = append(files, f)
//         return nil
//     })
//     return files, err
// }

// // FilesHandler scans both network folders and returns JSON
// func filesHandler(w http.ResponseWriter, r *http.Request) {
//     networks := []string{"_5000_network", "_5002_network"}
//     var allFiles []File

//     for _, netFolder := range networks {
//         files, err := scanNetworkFolder(netFolder)
//         if err != nil {
//             log.Println("Error scanning network folder:", netFolder, err)
//             continue
//         }
//         allFiles = append(allFiles, files...)
//     }

//     w.Header().Set("Content-Type", "application/json")
//     if err := json.NewEncoder(w).Encode(allFiles); err != nil {
//         http.Error(w, "Failed to encode files list", http.StatusInternalServerError)
//     }
// }

// func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
//     // Get the relative path from query
//     relPath := r.URL.Query().Get("name") // React should send the "Name" field
//     if relPath == "" {
//         http.Error(w, "Missing file name", http.StatusBadRequest)
//         return
//     }

//     // List of networks to check
//     networks := []string{"_5000_network", "_5002_network"}

//     var file *os.File
//     var fi os.FileInfo
//     var err error

//     // Try each network folder
//     for _, netFolder := range networks {
//         fullPath := filepath.Join(netFolder, relPath)
//         file, err = os.Open(fullPath)
//         if err == nil {
//             fi, err = file.Stat()
//             if err != nil {
//                 file.Close()
//                 continue
//             }
//             // Found the file
//             defer file.Close()
//             // Set headers and stream
//             mimeType := detectMimeType(fi.Name())
//             if mimeType == "" {
//                 mimeType = "application/octet-stream"
//             }

//             w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(fi.Name()))
//             w.Header().Set("Content-Type", mimeType)
//             w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
//             if _, err := io.Copy(w, file); err != nil {
//                 log.Println("Error writing file to response:", err)
//             }
//             return
//         }
//     }

//     // File not found in any network
//     log.Println("File not found in any network:", relPath)
//     http.Error(w, "File not found", http.StatusNotFound)
// }





///all of download




type File struct {
    ID          string `json:"ID"`
    Name        string `json:"Name"`
    Size        int64  `json:"Size"`
    Type        string `json:"Type"`
    CourseCode  string `json:"CourseCode,omitempty"`
    CourseName  string `json:"CourseName,omitempty"`
    Description string `json:"description,omitempty"`
    UploadedAt  int64  `json:"uploadedAt"`
}

// Directories for P2P networks
var networkDirs = []string{
    "_5000_network",
    "_5002_network",
}

// Map to keep track of file IDs to paths
var fileIndex = map[string]string{}

// helper to detect mime type
func detectMimeType(filename string) string {
    ext := strings.ToLower(filepath.Ext(filename))
    if ext != "" {
        mimeType := mime.TypeByExtension(ext)
        if mimeType != "" {
            return mimeType
        }
    }
    return "application/octet-stream"
}

// scan directories recursively and build index
func indexFiles() ([]File, error) {
    var allFiles []File
    fileIndex = map[string]string{} // reset

    for _, netDir := range networkDirs {
        root := filepath.Join(".", netDir)
        err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return nil
            }
            if !info.IsDir() {
                id := uuid.New().String()
                fileIndex[id] = path
                file := File{
                    ID:         id,
                    Name:       info.Name(),
                    Size:       info.Size(),
                    Type:       detectMimeType(info.Name()),
                    UploadedAt: info.ModTime().Unix(),
                }
                allFiles = append(allFiles, file)
            }
            return nil
        })
        if err != nil {
            return nil, err
        }
    }

    return allFiles, nil
}

// HTTP handler: list all files
func listFilesHandler(w http.ResponseWriter, r *http.Request) {
    files, err := indexFiles()
    if err != nil {
        http.Error(w, "Failed to scan files", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(files)
}

// HTTP handler: download by ID
func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "Missing file ID", http.StatusBadRequest)
        return
    }

    path, ok := fileIndex[id]
    if !ok {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    file, err := os.Open(path)
    if err != nil {
        http.Error(w, "Cannot open file", http.StatusInternalServerError)
        return
    }
    defer file.Close()

    fi, _ := file.Stat()
    mimeType := detectMimeType(fi.Name())
    w.Header().Set("Content-Disposition", "attachment; filename="+fi.Name())
    w.Header().Set("Content-Type", mimeType)
    w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))

    io.Copy(w, file)
}



//stop

// type File struct {
//     ID         string `json:"ID"`
//     Name       string `json:"Name"`
//     Size       int64  `json:"Size"`
//     Type       string `json:"Type"`
//     Path       string `json:"Path"`
//     UploadedAt int64  `json:"uploadedAt"`
// }

// var storageDirs = []string{"_5000_network", "_50002_network"}

// // Detect MIME by reading file header bytes
// func detectMimeType(path string) string {
//     f, err := os.Open(path)
//     if err != nil {
//         return "application/octet-stream"
//     }
//     defer f.Close()

//     header := make([]byte, 512)
//     n, _ := f.Read(header)
//     if n == 0 {
//         return "application/octet-stream"
//     }

//     mimeType := http.DetectContentType(header)
//     return mimeType
// }

// // Recursively scan P2P storage and extract file objects
// func FilesHandler(w http.ResponseWriter, r *http.Request) {
//     var files []File

//     for _, root := range storageDirs {
//         entries, err := os.ReadDir(root)
//         if err != nil {
//             continue
//         }

//         for _, e := range entries {
//             if e.IsDir() {
//                 continue
//             }

//             fi, err := e.Info()
//             if err != nil {
//                 continue
//             }

//             // skip tiny marker files
//             if fi.Size() < 5 {
//                 continue
//             }

//             files = append(files, File{
//                 ID:         e.Name(),                     // REAL CAS HASH
//                 Name:       e.Name(),                     // SAME AS ID
//                 Path:       filepath.Join(root, e.Name()),
//                 Size:       fi.Size(),
//                 Type:       detectMimeType(filepath.Join(root, e.Name())),
//                 UploadedAt: fi.ModTime().Unix(),
//             })
//         }
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(files)
// }


// // Download handler
// func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
//     name := r.URL.Query().Get("name")
//     if name == "" {
//         http.Error(w, "Missing CAS filename", http.StatusBadRequest)
//         return
//     }

//     // Search for CAS file in both folders
//     var fullPath string
//     for _, root := range storageDirs {
//         err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
//             if err != nil {
//                 return nil
//             }
//             if !info.IsDir() && info.Name() == name {
//                 fullPath = path
//             }
//             return nil
//         })
//         if fullPath != "" {
//             break
//         }
//         _ = err
//     }

//     if fullPath == "" {
//         http.Error(w, "File not found", http.StatusNotFound)
//         return
//     }

//     f, err := os.Open(fullPath)
//     if err != nil {
//         http.Error(w, "Cannot open file", http.StatusInternalServerError)
//         return
//     }
//     defer f.Close()

//     mimeType := detectMimeType(fullPath)
//     w.Header().Set("Content-Type", mimeType)

//     // Use hash as filename; real filename is lost in CAS
//     w.Header().Set("Content-Disposition", "attachment; filename="+name)

//     io.Copy(w, f)
// }




func main() {
    s1 := makeServer(":5000")           // no bootstrap nodes
    s2 := makeServer(":50002", ":5000")  // bootstrap node s1


    // Start both P2P nodes
    go func() { log.Println("Starting s1..."); log.Fatal(s1.Start()) }()
    go func() { log.Println("Starting s2..."); log.Fatal(s2.Start()) }()

    time.Sleep(2 * time.Second) // wait for bootstrap

    // HTTP server for uploads
    http.HandleFunc("/upload", uploadHandler(s2))
     http.HandleFunc("/files", listFilesHandler)
    http.HandleFunc("/files/download", downloadFileHandler)

    fmt.Println("Upload server running on :8082")
    log.Fatal(http.ListenAndServe(":8080", enableCORS(http.DefaultServeMux)))
}


