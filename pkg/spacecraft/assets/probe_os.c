#include <video.c>
#include <vfs.c>

#define INBOX_BUFFER 0x8000

// TODO: modify to use interrupts

int main() {
    print("Voyager-1 OS starting...\n");
    
    // Initial picture
    // take_picture_and_send("Earth");
    
    int* msgrecv = find_peripheral("MSGRECV");
    if (msgrecv == 0) {
        print("MSGRECV not found!\n");
        return 1;
    }
    
    while (1) {
        print("id:");
        // print(msgrecv);
        // print("\n");
        
        // Check MSGRECV state (offset 0)
        if (*msgrecv == 1) {
            print("Message received!\n");
            
            // Read command from INBOX.MSG
            char* inbox_file = "INBOX.MSG";
            int size = vfs_size_calc(inbox_file);
            if (size > 0) {
                vfs_read(inbox_file, INBOX_BUFFER);
                // Null terminate the command string
                char* cmd = (char*)INBOX_BUFFER;
                cmd[size] = 0;
                
                print("Command: ");
                print(cmd);
                print("\n");
                
                if (strcmp(cmd, "TAKE_PICTURE") == 0) {
                    print("Taking picture...\n");
                    take_picture_and_send("Earth");
                    print("sent...\n");
                } else {
                    print("Unknown command\n");
                }
            }
            
            // ACK message (clears it from queue)
            *msgrecv = 1; 
        }
        
        // Small delay to prevent pegging the CPU
        for (int i = 0; i < 5000; i++) {
            asm("NOP");
        }
    }
    
    return 0;
}
