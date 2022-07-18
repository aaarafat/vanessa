echo "Installing netfilter-queue..."
sudo apt-get update -qq | sudo apt-get install -y libnetfilter-queue-dev
sudo apt-get install python-is-python3 -y
echo "Installing netfilter-queue...done"
echo "Installing socat..."
sudo apt update
sudo apt install socat
echo "Installing socat...done"
echo "Installing python dependencies..."
sudo python -m pip install -r scripts/requirements.txt
echo "Installing python dependencies...done"
echo "Installing go dependencies..."
go get ./...
echo "Installing go dependencies...done"
echo "Installing npm dependencies..."
npm install
echo "Installing npm dependencies...done"