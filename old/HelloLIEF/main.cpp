#include <iostream>
#include <memory>
#include <unistd.h>
#include <string>
#include <fstream>
#include <LIEF/ELF.hpp>
#include <sys/types.h>
#include <sys/stat.h>

using namespace LIEF::ELF;

int main(int argc, char **argv)
{
  if (argc < 2)
  {
    std::cerr << "Usage: " << argv[0] << " <binary> <args>" << std::endl;
    return 1;
  }
  std::unique_ptr<LIEF::ELF::Binary> binary{LIEF::ELF::Parser::parse(argv[1])};

  Section section{".keyval"};
  section.type(ELF_SECTION_TYPES::SHT_PROGBITS);
  section.alignment(32);
  section.virtual_address(0x10000000);
  section.offset(0x10000000);
  section.add(ELF_SECTION_FLAGS::SHF_ALLOC);
  section.add(ELF_SECTION_FLAGS::SHF_WRITE);
  section.add(ELF_SECTION_FLAGS::SHF_EXECINSTR);
  std::vector<uint8_t> data(4 * 1024 * 1024);
  section.content(std::move(data));
  auto added_section = binary->add(section, true);
  for (Segment &segment : binary->segments())
  {
    if (segment.virtual_address() == added_section.virtual_address())
    {
      segment.alignment(0x1000);
      // segment.file_offset(0x10000000);
    }
  }

  std::string origName = std::string(argv[1]);
  std::string base_filename = origName.substr(origName.find_last_of("/\\") + 1);
  std::string newName = std::string("/kv-tmpfs/") + base_filename;
  char *newNameStr = (char *)newName.c_str();
  std::fstream fs;
  fs.open(newNameStr, std::ios::out);
  fs.close();
  binary->write(newNameStr);
  argv[1] = newNameStr;
  chmod(newNameStr, S_IRWXU);
  std::cout << "Created new file: " << newNameStr << std::endl;
  return execvp(argv[1], argv + 1);
}
