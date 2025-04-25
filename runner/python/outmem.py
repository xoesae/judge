import signal, os

os.kill(os.getpid(), signal.SIGSEGV)