# NASA HDF Data Processor

A batch processor built in Go for converting HDF (Hierarchical Data Format) files to JSON format. This tool is designed to process NASA Earth science data files and convert them to JSON files for further processing or analysis.

## Features

- 🔄 Convert HDF files (.h5, .hdf5, .he5, .nc) to JSON
- ⏰ Scheduled batch processing with configurable intervals
- 📊 Metadata extraction and data filtering
- 🔍 Automatic file discovery and processing
- 🐳 Docker containerization support
- 🚀 High-performance Go implementation
- 📝 Comprehensive logging and processing reports
- 🔄 Continuous monitoring with duplicate detection

## Usage Modes

### Continuous Processing
The processor runs continuously, checking for new HDF files at regular intervals and converting them to JSON format.


## Command Line Options

### File Processing Options
- `-input-dir` - Directory containing HDF files to process
- `-output-dir` - Directory to save converted JSON files
- `-interval` - Processing interval for continuous mode
- `-once` - Run once and exit (single-shot mode)

### Configuration Options
- `-log-level` - Logging level (debug, info, warn, error)
- `-format` - Output format (currently only JSON supported)

## Installation

### Prerequisites
- Go 1.21 or higher
- HDF5 library installed on your system

### From Source
```bash
git clone <repository-url>
cd nasa-hdf-processor
go mod download
go build -o nasa-hdf-processor .
```

### Using Docker
```bash
docker build -t nasa-hdf-processor .
docker run -p 8080:8080 -v /path/to/hdf/files:/app/data nasa-hdf-processor
```

### Using Docker Compose
```bash
# Place your HDF files in ./data directory
docker-compose up -d
```

## Usage

### Command Line Options
```bash
./nasa-hdf-processor -help
```

Options:
- `-port string` - Port to run server on (default: 8080)
- `-data-dir string` - Directory containing HDF files (default: ./data)
- `-log-level string` - Log level: debug, info, warn, error (default: info)

### Environment Variables
- `PORT` - Override server port
- `DATA_DIR` - Override data directory
- `GIN_MODE=release` - Run in production mode

### Examples

#### Basic Usage
```bash
./nasa-hdf-processor -input-dir /path/to/hdf/files -output-dir /path/to/json/output
```

#### Single Run Mode
```bash
./nasa-hdf-processor -once -input-dir ./hdf_files -output-dir ./json_output
```

#### Using Environment Variables
```bash
INPUT_DIR=/data/hdf OUTPUT_DIR=/data/json ./nasa-hdf-processor
```

#### Docker with Custom Directories
```bash
docker run -d \
  -v /path/to/hdf/files:/app/input:ro \
  -v /path/to/json/output:/app/output \
  -e PROCESS_INTERVAL=10m \
  nasa-hdf-processor
```

## Processing Examples

### Continuous Processing (Default)
```bash
# Process files every 5 minutes
./nasa-hdf-processor -input-dir ./input -output-dir ./output

# Process files every hour
./nasa-hdf-processor -input-dir ./input -output-dir ./output -interval 1h
```

### Single Run Processing
```bash
# Process all files once and exit
./nasa-hdf-processor -once -input-dir ./input -output-dir ./output
```

### Docker Compose
```bash
# Place HDF files in ./input directory
# JSON files will be generated in ./output directory
docker-compose up -d
```

## Output Files

The processor generates JSON files with timestamps to avoid conflicts:

```
output/
├── sample_20231201_143022.json
├── data_20231201_143025.json
├── processing_report_20231201_143030.json
└── ...
```

### Processing Reports
After each processing cycle, a report file is generated with statistics:

```json
{
  "timestamp": "2023-12-01T14:30:30Z",
  "total_files": 5,
  "success_count": 4,
  "failed_count": 1,
  "success_rate": 80.0,
  "total_size_mb": 125.5,
  "total_duration": "2m15s"
}
```

## Data Structure

The converted JSON structure includes:

```json
{
  "file_name": "sample.h5",
  "groups": {
    "group_name": {
      "datasets": { ... },
      "groups": { ... },
      "attributes": { ... }
    }
  },
  "datasets": {
    "dataset_name": {
      "rank": 2,
      "dimensions": [100, 200],
      "data_type": 1,
      "data": [...],
      "attributes": { ... }
    }
  },
  "attributes": {
    "global_attr": "value"
  }
}
```

## Performance Considerations

- Large datasets (>100MB) are skipped to prevent memory issues
- Files are processed concurrently (max 3 workers) to improve throughput
- Duplicate detection prevents reprocessing the same files within 1 hour
- Processing reports provide detailed statistics and performance metrics
- Output files include timestamps to avoid conflicts

## Supported File Formats

- HDF5 (.h5, .hdf5)
- HDF-EOS (.he5)
- NetCDF (.nc) - Note: Limited support, may require additional libraries

## Development

### Building
```bash
go build -o nasa-hdf-processor .
```

### Running Tests
```bash
go test ./...
```

### Code Formatting
```bash
go fmt ./...
```

### Linting
```bash
golangci-lint run
```

## Docker Development

### Build and Run
```bash
docker build -t nasa-hdf-processor .
docker run -v $(pwd)/input:/app/input -v $(pwd)/output:/app/output nasa-hdf-processor
```

### Development with Docker Compose
```bash
docker-compose up --build
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions:
1. Check the documentation
2. Search existing issues
3. Create a new issue with detailed information

## Roadmap

- [ ] Support for additional output formats (CSV, Parquet)
- [ ] Configurable file size limits and processing rules
- [ ] Integration with cloud storage (S3, GCS, Azure)
- [ ] Advanced filtering and data transformation options
- [ ] Web dashboard for monitoring processing status
- [ ] Integration with NASA Earthdata APIs
- [ ] Support for compressed HDF files
- [ ] Parallel processing optimization for large datasets
