class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage "https://github.com/tolitius/bb8"
  url "https://github.com/tolitius/bb8/releases/download/v0.1.11/bb8_0.1.11_macOS_amd64.tar.gz"
  version "0.1.11"
  sha256 "548963b58e2022a614db043d406b28e400af1c6654cbba78e1377cc43113d0da"

  def install
    bin.install "bb"
  end

  test do
    bb version
  end
end
