echo "Building Car Application....."
go build -o dist/apps/car apps/car/main.go
echo "Building Network....."
go build -o dist/apps/network apps/network/main.go
echo "Building Switch....."
go build -o apps/scripts/switch apps/scripts/switch.go
echo "Done Building!!!!"