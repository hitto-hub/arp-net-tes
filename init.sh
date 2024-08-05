#!/bin/bash

# Create a Python virtual environment
python3 -m venv venv

# Activate the virtual environment
source venv/bin/activate

# Install required packages
pip install -r requirements.txt

# Deactivate the virtual environment
deactivate

# Done
echo "Setup complete. To activate the virtual environment, run 'source venv/bin/activate'."