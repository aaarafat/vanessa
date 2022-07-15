echo "Installing netfilter-queue..."
sudo apt-get update -qq | sudo apt-get install -y libnetfilter-queue-dev
echo "Installing netfilter-queue...done"
echo "Installing socat..."
sudo apt update
sudo apt install socat
echo "Installing socat...done"
echo "Installing go dependencies..."
go get ./...
echo "Installing go dependencies...done"
echo "Installing npm dependencies..."
npm install
echo "Installing npm dependencies...done"