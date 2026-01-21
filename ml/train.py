import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import DataLoader
from dataset import SnakeOfflineRLDataset
from model import SnakePolicyNet
import os
import copy
import logging

# --- Monkey Patch for Python 3.14 / onnxscript compatibility ---
# The onnxscript library has a strict type check that fails with Python 3.14's typing system.
# We patch the constructor to ignore this specific error since the object is initialized correctly anyway.
try:
    import onnxscript.values
    _orig_init = onnxscript.values.AttrRef.__init__

    def _patched_init(self, attr_name, typeinfo, info):
        try:
            _orig_init(self, attr_name, typeinfo, info)
        except TypeError as e:
            if "Expecting a type not" in str(e):
                pass # Ignore the type check failure
            else:
                raise e

    onnxscript.values.AttrRef.__init__ = _patched_init
    print("🔧 Applied onnxscript Python 3.14 compatibility patch.")
except ImportError:
    pass
# ---------------------------------------------------------------

def train_rl():
    # 1. Config
    RECORDS_DIR = "../records"
    BATCH_SIZE = 64
    LR = 0.0005
    GAMMA = 0.95      # Discount factor for future rewards
    EPOCHS = 15
    TARGET_UPDATE = 2 # Update target network every N epochs
    DEVICE = torch.device("cuda" if torch.cuda.is_available() else "cpu")
    
    print(f"🚀 Starting Offline RL Training on {DEVICE}...")

    # 2. Data
    dataset = SnakeOfflineRLDataset(RECORDS_DIR)
    if len(dataset) == 0:
        print("❌ No data found.")
        return
    loader = DataLoader(dataset, batch_size=BATCH_SIZE, shuffle=True)

    # 3. Networks
    # Q-Network: the one we train
    q_net = SnakePolicyNet().to(DEVICE)
    # Target-Network: stable reference for Bellman targets
    target_net = copy.deepcopy(q_net).to(DEVICE)
    target_net.eval()

    optimizer = optim.Adam(q_net.parameters(), lr=LR)
    criterion = nn.MSELoss() # Q-learning uses MSE on TD-error

    # 4. Training Loop
    q_net.train()
    for epoch in range(EPOCHS):
        total_loss = 0
        
        for batch in loader:
            s = batch['s'].to(DEVICE)
            a = batch['a'].to(DEVICE)
            r = batch['r'].to(DEVICE)
            s_next = batch['s_next'].to(DEVICE)
            done = batch['done'].to(DEVICE)

            # --- Q-Learning Magic ---
            
            # 1. Get current Q-values for the action taken: Q(s, a)
            current_q = q_net(s).gather(1, a.unsqueeze(1)).squeeze(1)
            
            # 2. Get max Q-value for next state: max Q(s', a') from Target Net
            with torch.no_grad():
                next_q_max = target_net(s_next).max(1)[0]
                # Target = Reward + Gamma * MaxQ(next) * (1 - Done)
                target_q = r + (GAMMA * next_q_max * (1 - done))

            # 3. Compute Loss
            loss = criterion(current_q, target_q)
            
            # 4. Optimize
            optimizer.zero_grad()
            loss.backward()
            # Clip gradients to prevent exploding gradients
            torch.nn.utils.clip_grad_norm_(q_net.parameters(), 1.0)
            optimizer.step()
            
            total_loss += loss.item()

        # Update Target Network
        if epoch % TARGET_UPDATE == 0:
            target_net.load_state_dict(q_net.state_dict())

        avg_loss = total_loss / len(loader)
        print(f"📅 Epoch {epoch+1}/{EPOCHS} | TD-Loss: {avg_loss:.6f}")

    # 5. Save & Export
    os.makedirs("checkpoints", exist_ok=True)
    # Save PTH
    torch.save(q_net.state_dict(), "checkpoints/snake_rl_v1.pth")
    
    # Export ONNX
    q_net.eval()
    dummy_input = torch.randn(1, 6, 25, 25).to(DEVICE)
    torch.onnx.export(q_net, dummy_input, "checkpoints/snake_policy.onnx",
                      opset_version=18,  # Stable opset for modern PyTorch
                      input_names=['input'], 
                      output_names=['output'])
        
    print("✨ RL Model Trained and Exported to ONNX!")

if __name__ == "__main__":
    train_rl()
