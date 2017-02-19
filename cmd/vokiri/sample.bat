@echo off
setlocal ENABLEEXTENSIONS
rem %~dp0 �܂ł��A���̃t�@�C��������t�H���_�{��
set DEST=%~dp0kiritan_
rem set DEBUG=--debug

pushd %~dp0

if "%1"=="help" (
    call :help
) else if "%1"=="gen" (
    call :gen
) else if "%1"=="clean" (
    call :clean
) else  (
    call :sample
)

popd
endlocal
echo on

@exit /b 0


:sample

vokiri --volume=1 --speed=0.95 --pitch=1 --emph=1.2 --persist

vokiri %DEBUG% ������%date:~0,4%�N%date:~5,2%��%date:~8,2%���ł��B
vokiri %DEBUG% VOKIRI�ɋ����������Ă����������肪�Ƃ��������܂��B
vokiri %DEBUG% VOKIRI�́A�R�}���h���C������VOICEROID�𑀍삷�邽�߂̃\�t�g�E�F�A�ł��B
vokiri %DEBUG% ���Ƃ��팸������A�ق��̃\�t�g����̋��n�����ɂȂ����肷�邱�Ƃ����҂��Ă��܂��B

vokiri %DEBUG% ���āAVOKIRI�ł́A�኱�ނ����VOICEROID�̃E�B���h�E�𑀍삵�Ă��邱�Ƃ����聃�������āA���ꂵ�����@���Ȃ���ł�����b����
vokiri %DEBUG% �H�ɓ������~�܂��Ă��܂����Ƃ�����܂��B
vokiri %DEBUG% ���̏ꍇ�́A������񓯂����Ƃ����s���Ă݂Ă��������B���������悤�ɂȂ�Ǝv���܂��B
vokiri %DEBUG% ����ł����߂Ȃ�A�R�}���h���C���̃I�v�V�����Ƃ��āA����--debug�b�͂��ӂ�͂��ӂ�ł΂����������w�肵�āA�~�܂���������҂ɋ����Ă����Ă��������B

vokiri %DEBUG% ���A���ƁA���k���肽��ȊO��VOICEROID�ł̓���󋵂ɂ��Ă��A�����Ă���������ƍK���ł��B
vokiri %DEBUG% ����--exe�b�͂��ӂ�͂��ӂ� ���[�����������[�����A�ŁA���삳������VOICEROID�̎��s�t�@�C�����A
vokiri %DEBUG% ����--title�b�͂��ӂ�͂��ӂ� �����Ƃ遄���A�ŁA����VOICEROID�̃E�B���h�E�^�C�g���i"VOICEROID �՗t��"���j���w�肵�Ă݂Ă��������B

vokiri %DEBUG% �ł͂ł́B

ping -w 0 -n 2 0.0.0.0>nul
vokiri %DEBUG% --volume=0.5 --pitch=2 --emphasis=2 ���ݶܲ�ԯ��
vokiri %DEBUG% --volume=0.5 --pitch=2 --emphasis=2 ��Ű�@

exit /b 0


:clean

del kiritan*

exit /b 0


:gen

echo %TIME%

rem
rem �Љ������̉����t�@�C���쐬�B
rem

vokiri --volume=1 --speed=0.95 --pitch=1 --emph=1.2 --persist

rem
rem ���A
rem
vokiri %DEBUG% --record-once="%DEST%0001.wav" "����ɂ��́A�i�i���k���肽��b(spd:0.8)�g�[�z�N(pit:1.20)�L!���^��$1_1�j�j�i�i�ł��b(pit:1.2)�f�b�XD�j�j�B "
vokiri %DEBUG% --record-once="%DEST%0002.wav" ���̓���ł́AVOKIRI�Ƃ���VOICEROID�x���\�t�g�̏Љ�ƁA�����܂ł̊ȒP�Ȑ��������܂��B
vokiri %DEBUG% --record-once="%DEST%0003.wav" VOKIRI�́AVOICEROID�����g���̕������̃\�t�g�ŁA��ɋ�����@�\�������Ă��܂��B
vokiri %DEBUG% --record-once="%DEST%0004.wav" ���Ȃ݂ɁA���݂͓��k���肽��̑����ʂɍ��킹�č�肱��ł��܂����A�����I�ɂ͑���VOICEROID�ɂ��Ή�����c��������܂���B
rem �R�}���h���C���ł̓ǂݏグ
vokiri %DEBUG% --record-once="%DEST%0005.wav" �@�\���̂����B�R�}���h���C����VOICEROID�{ ���k���肽��ɓǂݏグ������
vokiri %DEBUG% --record-once="%DEST%0006.wav" ��{�@�\�ł��B
rem �C���g�l�[�V�������̒���
vokiri %DEBUG% --record-once="%DEST%0007.wav" �@�\���̂ɁB�C���g�l�[�V�������̒������R�}���h���C���Ŏw�肷��
vokiri %DEBUG% --record-once="%DEST%0008.wav" �ʏ�̃R�}���h���C���c�[�����ł͂��܂�������Ă��Ȃ��A�R�}���h���C������̒������s���܂��B
rem �����t�@�C���ւ̕ۑ�
vokiri %DEBUG% --record-once="%DEST%0009.wav" �@�\���̂���B�ǂݏグ�����e�������t�@�C���Ƃ��ĕۑ�����
vokiri %DEBUG% --record-once="%DEST%0010.wav" �����̃c�[���Ŏ����ł��Ă��邱�Ƃł����A���RVOKIRI�ł��s���܂��B
vokiri %DEBUG% --record-once="%DEST%0011.wav" ���̓�����A���̕��@�� �������k���肽��b�킽�������ɓǂݏグ�����Ă��܂��B
rem ���ύX�̉����t�@�C���ۑ��̃X�L�b�v
vokiri %DEBUG% --record-once="%DEST%0012.wav" �@�\���̂��B���ύX�̉����t�@�C���ۑ��̃X�L�b�v
vokiri %DEBUG% --record-once="%DEST%0013.wav" �����̃R�}���h���C���ŉ����t�@�C������C�ɏ����o�������ꍇ�A�o�͂ɂ͎��Ԃ�������܂��B
vokiri %DEBUG% --record-once="%DEST%0014.wav" ���̋@�\���g���ƁA�������񉹐��t�@�C�����o�͂��������ňꕔ��ύX�����ꍇ�̏o�̓X�s�[�h��啝�ɉ��P���܂��B
rem
vokiri %DEBUG% --record-once="%DEST%0015.wav" VOKIRI�͂����������@�\�������Ă��܂��B
vokiri %DEBUG% --record-once="%DEST%0016.wav" "�i�i���p���Ă��������ˁb(pit:0.95)�J�c���E �V�e (pit:1.1)�N^�_�T!�C(spd:2)(emph:0.5)�l^�[<R>�j�j�B"

rem
rem �_�E�����[�h
rem
vokiri %DEBUG% --record-once="%DEST%0100.wav" �_�E�����[�h�̂��[���[���[
vokiri %DEBUG% --record-once="%DEST%0101.wav" �܂��́A���茳�ɁA�킽�������k���肽������p�ӂ��������B
vokiri %DEBUG% --record-once="%DEST%0102.wav" --emph=2 �����A�����ꖳ���悤�ł�����A������@�ɂ��w�����������I
vokiri %DEBUG% --record-once="%DEST%0103.wav" --speed=0.8 ���ق�A
vokiri %DEBUG% --record-once="%DEST%0104.wav" �����āA����VOKIRI���_�E�����[�h���Ă��������B
vokiri %DEBUG% --record-once="%DEST%0105.wav" �_�E�����[�h�ꏊ�́A����̃R�����g���ɋL�ڂ��Ă���܂��B
vokiri %DEBUG% --record-once="%DEST%0106.wav" �_�E�����[�h������AZIP�t�@�C����W�J���܂��B

rem
rem �g���Ă݂�
rem
vokiri %DEBUG% --record-once="%DEST%0200.wav" �g���Ă݂܂��傤�I
vokiri %DEBUG% --record-once="%DEST%0201.wav" �_�E�����[�h���ł�����A���ɂ��遃��sample�b�T���v�������Ƃ����t�@�C�������s���Ă݂Ă��������B
vokiri %DEBUG% --record-once="%DEST%0202.wav" ���A���k���肽�񂪉����ǂݏグ�͂��߂܂��B
vokiri %DEBUG% --record-once="%DEST%0203.wav" ���̃t�@�C���ɂ́AVOKIRI���g���Ď��ɓǂݏグ������A�Ƃ����������L�q����Ă��܂��B
vokiri %DEBUG% --record-once="%DEST%0204.wav" ���̂ق��ɂ��A���̓���Ɏg�����߂̉����t�@�C���o�͂̃R�}���h���l�܂��Ă��܂��B
vokiri %DEBUG% --record-once="%DEST%0205.wav" ���l�Ɂ��������b�ǂ����񁄁�����Ă��� ����README�b��[�ǂ݁[���� �Ƃ��ǂ��A�Q�l�ɂȂ邩�Ǝv���܂��B

rem
rem ��������
rem
vokiri %DEBUG% --record-once="%DEST%0300.wav" ����ȂƂ���ł��I
vokiri %DEBUG% --record-once="%DEST%0301.wav" �����VOKIRI�̏Љ�Ɛ������I���܂��B
vokiri %DEBUG% --record-once="%DEST%0302.wav" ���������ȒP�ȏЉ�ɂȂ�܂������A���e�͂��������B
vokiri %DEBUG% --record-once="%DEST%0303.wav" ���������b�ǂ����񁄁�����Ă��� ����README�b��[�ǂ݁[���� �ɂ́A���ڍׂȋ@�\�̎g�������L�ڂ���Ă��܂��B
vokiri %DEBUG% --record-once="%DEST%0304.wav" �ق��̃R�}���h���o�͂����e�L�X�g(�W���o��)���󂯎���ēǂݏグ�������i���Љ�Ă��܂��B
vokiri %DEBUG% --record-once="%DEST%0305.wav" ���p���Ă���������Ɗ������ł��B
vokiri %DEBUG% --record-once="%DEST%0306.wav" �ł͂ł́[�I

rem
rem ���܂�
rem
vokiri --record-once="%DEST%omake_kwaii.wav" --pitch=2 --emphasis=2 ���ݶܲ�ԯ���@
vokiri --record-once="%DEST%omake_seya.wav" --pitch=2 --emphasis=2 ��Ű�@

vokiri �����܁[�[���I

dir /b *.wav > kiritan.m3u

echo %TIME%

exit /b 0


:help

echo HELP:
echo   %~n0 gen: ����p��WAV�t�@�C�����o�͂��܂��B�e�L�X�g�t�@�C�����o�͂���ݒ�̏ꍇ�A�ύX�������̂������������܂��B
echo   %~n0 clean: �o�͂��ꂽ�t�@�C���ikiritan�`�j���폜���܂��B

exit /b 0

