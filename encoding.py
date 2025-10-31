import chardet

with open('./files/out/SEGMENTO.TXT', 'r', encoding='latin-1') as f:
    for line in f.readlines():
        print(line)


