import setuptools

with open('VERSION.txt', 'r') as f:
    version = f.read().strip()

setuptools.setup(
    name="saturn10-addons-smart-assembly",
    description="Meta package for oca-web Saturn addons",
    version=version,
    install_requires=[
        'saturn10-addons-sa-sa_base',
    ],
    classifiers=[
        'Programming Language :: Python',
        'Framework :: Saturn',
    ]
)