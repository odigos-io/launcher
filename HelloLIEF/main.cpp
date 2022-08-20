#include <iostream>
#include <memory>
#include <unistd.h>
#include <LIEF/ELF.hpp>

using namespace LIEF::ELF;

int main(int argc, char** argv) {
  if (argc < 2) {
    std::cerr << "Usage: " << argv[0] << " <binary> <args>" << std::endl;
    return 1;
  }

  std::unique_ptr<LIEF::ELF::Binary> binary{LIEF::ELF::Parser::parse(argv[1])};
  Section section{".keyval"};
  section.type(ELF_SECTION_TYPES::SHT_PROGBITS);
  section.alignment(0x1000);
  section.add(ELF_SECTION_FLAGS::SHF_ALLOC);
  section.add(ELF_SECTION_FLAGS::SHF_WRITE);
  std::vector<uint8_t> data(4096);
  section.content(std::move(data));
  binary->add(section, true);
  binary->write(argv[1]);
  return execvp(argv[1], argv + 1);
}
