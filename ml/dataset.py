import json
import glob
import os
import torch
import numpy as np
from torch.utils.data import Dataset

# Mapping direction vectors to class indices
DIR_MAP = {
    (0, -1): 0,
    (0, 1): 1,
    (-1, 0): 2,
    (1, 0): 3
}

class SnakeOfflineRLDataset(Dataset):
    def __init__(self, records_dir, grid_width=25, grid_height=25):
        self.files = glob.glob(os.path.join(records_dir, "*.jsonl"))
        self.transitions = []
        self.grid_width = grid_width
        self.grid_height = grid_height
        
        print(f"🔄 Parsing {len(self.files)} files for RL transitions...")
        self._load_data()
        print(f"✅ Loaded {len(self.transitions)} transitions (S, A, R, S', Done).")

    def _load_data(self):
        for fpath in self.files:
            game_steps = []
            with open(fpath, 'r') as f:
                for line in f:
                    try:
                        record = json.loads(line)
                        game_steps.append(record)
                    except:
                        continue
            
            # Form pairs (S, A, R, S', Done)
            for i in range(len(game_steps) - 1):
                curr = game_steps[i]
                nxt = game_steps[i+1]
                
                # Extract Action
                action_obj = curr['action']
                dir_obj = action_obj.get('dir', {})
                dx, dy = dir_obj.get('x', 0), dir_obj.get('y', 0)
                a_idx = DIR_MAP.get((dx, dy))
                
                if a_idx is not None:
                    self.transitions.append({
                        's': curr['state'],
                        'a': a_idx,
                        'r': curr.get('reward', 0.0),
                        's_next': nxt['state'],
                        'done': curr.get('done', False)
                    })

    def _state_to_tensor(self, state):
        grid = np.zeros((6, self.grid_height, self.grid_width), dtype=np.float32)
        
        def set_p(c, p_list):
            for p in p_list:
                # Handle both {x,y} and {pos:{x,y}} formats
                x, y = -1, -1
                if 'x' in p: x, y = p['x'], p['y']
                elif 'pos' in p: x, y = p['pos']['x'], p['pos']['y']
                
                if 0 <= x < self.grid_width and 0 <= y < self.grid_height:
                    grid[c, y, x] = 1.0

        set_p(0, (state.get('snake') or [])[:1]) # Player Head
        set_p(1, (state.get('snake') or [])[1:]) # Body
        set_p(2, (state.get('aiSnake') or [])[:1]) # Enemy Head
        set_p(3, (state.get('aiSnake') or [])[1:]) # Enemy Body
        set_p(4, state.get('foods') or [])
        
        # Hazards
        for obs in (state.get('obstacles') or []):
            set_p(5, obs.get('points') or [])
        set_p(5, state.get('fireballs') or [])
        
        return torch.from_numpy(grid)

    def __len__(self):
        return len(self.transitions)

    def __getitem__(self, idx):
        t = self.transitions[idx]
        s = self._state_to_tensor(t['s'])
        s_next = self._state_to_tensor(t['s_next'])
        
        return {
            's': s,
            'a': torch.tensor(t['a'], dtype=torch.long),
            'r': torch.tensor(t['r'], dtype=torch.float32),
            's_next': s_next,
            'done': torch.tensor(1.0 if t['done'] else 0.0, dtype=torch.float32)
        }
