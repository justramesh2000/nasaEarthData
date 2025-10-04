#!/bin/bash

# NASA HDF Data Processor - Test Setup Script

set -e

echo "🚀 NASA HDF Data Processor - Test Setup"
echo "======================================="

# Create directories
echo "Creating directories..."
mkdir -p input output logs

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Check if HDF5 is installed (basic check)
if ! pkg-config --exists hdf5 2>/dev/null; then
    echo "⚠️  HDF5 library not found via pkg-config. You may need to install HDF5 development libraries."
    echo "   On Ubuntu/Debian: sudo apt-get install libhdf5-dev"
    echo "   On macOS: brew install hdf5"
    echo "   On CentOS/RHEL: sudo yum install hdf5-devel"
else
    echo "✅ HDF5 library found: $(pkg-config --modversion hdf5)"
fi

# Install Go dependencies
echo "Installing Go dependencies..."
go mod download
go mod tidy

# Try to build the application
echo "Building application..."
if go build -o nasa-hdf-processor .; then
    echo "✅ Build successful"
else
    echo "❌ Build failed. Please check the error messages above."
    exit 1
fi

# Test the help command
echo "Testing application..."
if ./nasa-hdf-processor -help > /dev/null 2>&1; then
    echo "✅ Application runs correctly"
else
    echo "❌ Application failed to run"
    exit 1
fi

# Create a sample test file (if possible)
echo "Creating sample test structure..."
cat > input/README.txt << EOF
This is the input directory for HDF files.

Place your HDF files (.h5, .hdf5, .he5, .nc) in this directory.
The processor will automatically detect and convert them to JSON format.

Example files:
- sample.h5
- data.hdf5
- satellite.he5
- measurement.nc
EOF

cat > output/README.txt << EOF
This is the output directory for converted JSON files.

Converted files will be saved here with timestamps:
- sample_20231201_143022.json
- data_20231201_143025.json
- processing_report_20231201_143030.json
EOF

echo ""
echo "🎉 Setup completed successfully!"
echo ""
echo "Next steps:"
echo "1. Place HDF files in the 'input' directory"
echo "2. Run: ./nasa-hdf-processor -input-dir input -output-dir output"
echo "3. Or run once: ./nasa-hdf-processor -once -input-dir input -output-dir output"
echo ""
echo "For continuous processing:"
echo "  ./nasa-hdf-processor -input-dir input -output-dir output -interval 5m"
echo ""
echo "For Docker:"
echo "  make docker-run"
echo ""
echo "For help:"
echo "  ./nasa-hdf-processor -help"
echo ""
