#include <video.c>
#include <vfs.c>
#include <stdio.c>

#define INBOX_BUFFER 0x8000

int* INT_MASK = 0xFF09;
int* MMIO_SLOT_BASE = 0xFE00;
char* COMMAND_TAKE_PICTURE = "TAKE_PICTURE";

void isr() {
    int pending = *INT_MASK;

    int* slot_ptr = find_peripheral("MSGRECV");
    if (slot_ptr == 0) {
        print("Error: Message Receiver Peripheral not found!\n");
        return;
    }

    int address = (int)slot_ptr;
    int offset = address - 0xFE00;
    int RECV_SLOT = offset / 16;

    if (RECV_SLOT != -1) { // this is the slot where the message reciever is attatched
        int mask = 1;
        for (int i = 0; i < RECV_SLOT; i++) {
            mask = mask * 2;
        }

        if ((pending & mask) != 0) {
            print("new message");

            char buffer[256];
            char sender_buffer[256];
            char* filename = "INBOX.MSG";
            char* sender_filename = "SENDER.MSG";

            int size = vfs_size_calc((int*)filename);
            int sender_size = vfs_size_calc((int*)sender_filename);

            if (size >= 0 && sender_size >= 0) {
                if (size < 255 && sender_size < 255) {
                    int err_sender = vfs_read((int*)sender_filename, (int*)sender_buffer);
                    int err_msg = vfs_read((int*)filename, (int*)buffer);
                    
                    if (err_sender == 0 && err_msg == 0) {
                        sender_buffer[sender_size] = 0;
                        buffer[size] = 0;
                        
                        print("Message Received from ");
                        print(sender_buffer);
                        print(": ");
                        print(buffer);
                        print("\n");

                        if (strcmp(COMMAND_TAKE_PICTURE, (char*)buffer) == 0  ){
                            take_picture_and_send(sender_buffer);
                        } else {
                            print("Unknown message:");
                            print(buffer);
                        }
                    } else {
                        print("Error reading messages. Sender err: ");
                        print_int(err_sender);
                        print(", Msg err: ");
                        print_int(err_msg);
                        print("\n");
                    }
                } else {
                    print("Error: Message or Sender too large\n");
                }
                vfs_delete((int*)filename);
                vfs_delete((int*)sender_filename);
            } else {
                print("Error: INBOX.MSG or SENDER.MSG not found or invalid\n");
            }

            int* slot_addr = 0xFE00 + (RECV_SLOT * 16);
            *slot_addr = 1;

            // Only clear our interrupt
            *INT_MASK = mask;
        }
    }
}

int main() {
    print("Voyager-1 OS starting...\n");
    
    enable_interrupts();
    print("Interrupts enabled. Waiting for messages...\n");

    while (1) {
        wait_for_interrupt();
    }
    
    return 0;
}
