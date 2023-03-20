#!/bin/bash

# Set variables
SERVICE_NAME="<ServiceName>"
EXECUTABLE_PATH="<ExecutablePath>"
WORKING_DIRECTORY="<WorkingDirectory>"

# Create systemd service file
cat > /etc/systemd/system/$SERVICE_NAME.service <<EOL
[Unit]
Description=$SERVICE_NAME

[Service]
ExecStart=$EXECUTABLE_PATH
WorkingDirectory=$WORKING_DIRECTORY

[Install]
WantedBy=multi-user.target
EOL

# Reload systemd daemon and enable service to start at boot time.
systemctl daemon-reload && systemctl enable $SERVICE_NAME.service && systemctl start $SERVICE_NAME