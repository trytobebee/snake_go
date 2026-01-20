import torch
import torch.nn as nn
import torch.nn.functional as F

class SnakePolicyNet(nn.Module):
    def __init__(self, input_channels=6, num_actions=4):
        super(SnakePolicyNet, self).__init__()
        
        # Simple CNN Architecture
        # Input: (B, 6, 25, 25)
        
        # Layer 1
        self.conv1 = nn.Conv2d(input_channels, 32, kernel_size=3, padding=1)
        self.bn1 = nn.BatchNorm2d(32)
        
        # Layer 2
        self.conv2 = nn.Conv2d(32, 64, kernel_size=3, padding=1)
        self.bn2 = nn.BatchNorm2d(64)
        
        # Layer 3
        self.conv3 = nn.Conv2d(64, 64, kernel_size=3, padding=1)
        self.bn3 = nn.BatchNorm2d(64)
        
        # Fully Connected
        # 25x25 grid after pooling/strides? 
        # Actually keeping spatial dim constant is better for board games usually,
        # but let's do a simple flatten here.
        self.fc1 = nn.Linear(64 * 25 * 25, 512)
        self.fc2 = nn.Linear(512, num_actions) # 4 Directions
        
        # Auxiliary heads for Fire/Boost could be added here later
        
    def forward(self, x):
        # x: (B, 6, 25, 25)
        
        x = F.relu(self.bn1(self.conv1(x)))
        x = F.relu(self.bn2(self.conv2(x)))
        x = F.relu(self.bn3(self.conv3(x)))
        
        x = x.reshape(x.size(0), -1) # Flatten (B, 64*25*25)
        
        x = F.relu(self.fc1(x))
        logits = self.fc2(x)
        
        return logits
