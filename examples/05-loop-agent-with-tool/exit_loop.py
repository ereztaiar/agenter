def exit_loop():
              """Call this function ONLY when the critique is 'APPROVED', indicating the story is finished and no more changes are needed."""
              print('{"status": "approved", "message": "Story approved. Exiting refinement loop."}')
      
if __name__ != "__main__":
    exit_loop()
    