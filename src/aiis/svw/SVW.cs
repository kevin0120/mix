using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using core;

namespace svw
{
    public class SVW : IPlugin
    {
        public void Init()
        {
            _fis_client = new FisClinet();
            _fis_client.Init();

        }

        public void Release()
        {
            Console.WriteLine("test_plugin release");
        }

        public void Run()
        {
            Console.WriteLine("test_plugin run");
        }

        private FisClinet _fis_client;
    }
}
