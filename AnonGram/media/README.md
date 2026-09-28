# FFmpeg media worker

Оригиналы файлов хранятся в S3/MinIO. Задания публикуются в Redis Stream `anongram:media`.
FFmpeg используется для thumbnail, preview, audio normalization и transcoding.
Пример: `ffmpeg -i input.mp4 -vf "scale=720:-2" -t 15 preview.mp4`.
