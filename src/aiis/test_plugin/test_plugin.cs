using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using core;

namespace test_plugin
{
    public class TestPlugin : IPlugin
    {
        public void Init()
        {
            Console.WriteLine("test_plugin init");
        }

        public void Release()
        {
            Console.WriteLine("test_plugin release");
        }

        public void Run()
        {
            Console.WriteLine("test_plugin run");
        }
    }
}
