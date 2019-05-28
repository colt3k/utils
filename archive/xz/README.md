
github.com/ulikunitz/xz

While I share a lot of the observations (complex format, unclear padding rules, LZMA2 
format undocumented), I don't share all the conclusions particularly those about the 
severity of bit errors in the headers. 

If your archiving system doesn't protect you from bit errors then you shouldn't store 
compressed data on it. In practical terms a bit error in a compressed file will affect 
all data following it. The exception are parallel compressed files, where unavoidable
dictionary resets provide synchronization points. The parallel-compressed segments could 
be quite large (>= 8 MiB); compare that to an uncompressed UTF-8 file, where every code 
point is a synchronization point.

If xz is used for archiving purposes, I recommend to use it without any additional filters.
You should use the SHA256 checksum, since it also protect the length of the data stream. 
CRC-32 and CRC-64 lack this capability. This might contradict a statement from the article,
 but SHA256 embeds the length of the data stream into the data that is hashed. The fact 
 that the length field is not secured by a check sum, doesn't matter.

NOTE:
I don't support any additional filters, but the default is the CRC64 checksum. You would 
need a WriterConfig that sets the CheckSum field to SHA256.
