    
 export const findRoom = roomName => {
        for (let i = 0; i < this.rooms.length; i++) {
          if (this.rooms[i].name === roomName) {
              return this.rooms[i];
          }
        }
    },