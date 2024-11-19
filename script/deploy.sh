scp ~/codespeakss/lumipigeon/web/* 67_root:/var/www/html/;
ssh 67_root "pkill ap" ; scp ~/codespeakss/lumipigeon/out/ap 67_root:/root ; ssh 67_root " /root/ap"