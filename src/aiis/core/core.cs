using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace core
{
    public interface IPlugin
    {
        void Init();

        void Release();

        void Run();

    }
}
